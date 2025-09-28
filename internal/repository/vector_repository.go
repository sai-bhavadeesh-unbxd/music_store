package repository

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/redis/go-redis/v9"
)

type VectorRepository interface {
	UpsertSong(name string, embedding []float32) error
	UpsertUser(id string, embedding []float32) error
	SimilarSongsByVector(vec []float32, count int) ([]string, error)
}

type vectorRepository struct {
	rdb *redis.Client
}

const (
	songsSetKey    = "vs:songs"
	usersSetKey    = "vs:users"
	songsIdx       = "idx:songs"
	songHashPrefix = "vec:song:"
)

func NewVectorRepository(rdb *redis.Client) VectorRepository {
	return &vectorRepository{rdb: rdb}
}

func (v *vectorRepository) UpsertSong(name string, embedding []float32) error {
	// Ensure RediSearch index exists (HASH) before upsert
	_ = v.ensureHashIndex(len(embedding))
	// Upsert HASH for RediSearch vector indexing
	if err := v.rdb.HSet(context.Background(), songHashKey(name), "embedding", float32ToBytesLE(embedding), "name", name).Err(); err != nil {
		return err
	}
	return nil
}

func (v *vectorRepository) UpsertUser(id string, embedding []float32) error {
	values := float32To64(embedding)
	if _, err := v.rdb.VAdd(context.Background(), usersSetKey, id, &redis.VectorValues{Val: values}).Result(); err != nil {
		return err
	}
	return nil
}

func (v *vectorRepository) SimilarSongsByVector(vec []float32, count int) ([]string, error) {
	values := float32To64(vec)
	if res, err := v.rdb.VSimWithArgs(context.Background(), songsSetKey, &redis.VectorValues{Val: values}, &redis.VSimArgs{Count: int64(count)}).Result(); err == nil {
		return res, nil
	} else if !unknownVSim(err) {
		return nil, err
	}
	// RediSearch KNN on HASH fallback (requires FT.CREATE idx:songs ON HASH PREFIX 1 vec:song: ...)
	if err := v.ensureHashIndex(len(vec)); err != nil {
		// proceed to JSON fallback if index creation fails
	}
	if res, err := v.ftSearchTopK(vec, count); err == nil {
		if len(res) > 0 {
			return res, nil
		}
	}

	// Final fallback: cosine over JSON docs
	if res, err := v.topKFromJSON(vec, count); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

func float32To64(in []float32) []float64 {
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = float64(v)
	}
	return out
}

func float32ToBytesLE(in []float32) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, len(in)*4))
	for _, f := range in {
		_ = binary.Write(buf, binary.LittleEndian, f)
	}
	return buf.Bytes()
}

func (v *vectorRepository) ftSearchTopK(vec []float32, k int) ([]string, error) {
	blob := float32ToBytesLE(vec)
	// Build the KNN query as a single string so FT.SEARCH parses it correctly
	// Example: *=>[KNN 5 @embedding $vec AS score]
	query := fmt.Sprintf("*=>[KNN %d @embedding $vec AS score]", k)
	// Single attempt; if empty, caller will fallback to JSON
	args := []interface{}{songsIdx, query, "PARAMS", 2, "vec", blob, "WITHSCORES", "NOCONTENT", "SORTBY", "score", "ASC", "LIMIT", 0, k, "DIALECT", 2}
	res, err := v.rdb.Do(context.Background(), append([]interface{}{"FT.SEARCH"}, args...)...).Result()
	if err != nil {
		return nil, err
	}
	// Parse map-style response: { results: [ {id: ..., score: ..., values: [...]}, ...] }
	if m := normalizeAnyMap(res); m != nil {
		if raw, ok := m["results"]; ok {
			if list, ok := raw.([]interface{}); ok {
				out := make([]string, 0, k)
				for _, it := range list {
					mm := normalizeAnyMap(it)
					if mm == nil {
						continue
					}
					idStr := fmt.Sprintf("%v", mm["id"])
					if after, ok0 := strings.CutPrefix(idStr, songHashPrefix); ok0 {
						idStr = after
					}
					out = append(out, idStr)
					if len(out) >= k {
						break
					}
				}
				return out, nil
			}
		}
	}
	// Unknown/empty layout -> return empty to allow JSON fallback
	return []string{}, nil
}

// normalizeAnyMap converts map[string]interface{} or map[interface{}]interface{} to map[string]interface{}; otherwise returns nil
func normalizeAnyMap(in interface{}) map[string]interface{} {
	switch t := in.(type) {
	case map[string]interface{}:
		return t
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, v := range t {
			out[fmt.Sprintf("%v", k)] = v
		}
		return out
	default:
		return nil
	}
}

// JSON cosine fallback
type songDoc struct {
	Name      string    `json:"name"`
	Embedding []float32 `json:"embedding"`
}

func (v *vectorRepository) topKFromJSON(query []float32, k int) ([]string, error) {
	keys, err := v.rdb.Keys(context.Background(), "song:*").Result()
	if err != nil {
		return nil, err
	}
	type item struct {
		name  string
		score float64
	}
	items := make([]item, 0, len(keys))
	for _, key := range keys {
		// Try RedisJSON first
		s, err := v.rdb.Do(context.Background(), "JSON.GET", key).Text()
		if err != nil {
			// Fallback: plain GET if key is a simple string
			s, err = v.rdb.Get(context.Background(), key).Result()
			if err != nil {
				continue
			}
		}
		var doc songDoc
		if err := json.Unmarshal([]byte(s), &doc); err != nil {
			continue
		}
		if len(doc.Embedding) == 0 || len(query) == 0 {
			continue
		}
		items = append(items, item{name: doc.Name, score: cosine(query, doc.Embedding)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].score > items[j].score })
	if k > len(items) {
		k = len(items)
	}
	out := make([]string, 0, k)
	for i := 0; i < k; i++ {
		out = append(out, items[i].name)
	}
	return out, nil
}

func cosine(a []float32, b []float32) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var dot, na, nb float64
	for i := 0; i < n; i++ {
		fa := float64(a[i])
		fb := float64(b[i])
		dot += fa * fb
		na += fa * fa
		nb += fb * fb
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func unknownVSim(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "unknown command 'vsim'")
}

func songHashKey(name string) string { return songHashPrefix + name }

// ensureHashIndex creates a RediSearch index on HASH keys if missing
func (v *vectorRepository) ensureHashIndex(dim int) error {
	// If index exists, return
	if _, err := v.rdb.Do(context.Background(), "FT.INFO", songsIdx).Result(); err == nil {
		return nil
	}
	// Create index with explicit HNSW params and name field for debugging
	_, err := v.rdb.Do(
		context.Background(),
		"FT.CREATE", songsIdx,
		"ON", "HASH",
		"PREFIX", 1, songHashPrefix,
		"SCHEMA",
		"name", "TEXT", // not used in KNN, but helpful for inspection
		"embedding", "VECTOR", "HNSW", 10, // number of args below
		"TYPE", "FLOAT32",
		"DIM", dim,
		"DISTANCE_METRIC", "COSINE",
		"M", 16,
		"EF_CONSTRUCTION", 200,
	).Result()
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "index already exists") {
		return err
	}
	return nil
}
