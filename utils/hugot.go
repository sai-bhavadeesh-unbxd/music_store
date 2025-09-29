package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/knights-analytics/hugot"
)

// HugotConfig holds configuration for the embedding provider (e.g., model name, API key).
type HugotConfig struct {
	ModelName string
	Dim       int
	OnnxPath  string
	ModelsDir string
}

// GetDefaultHugotConfig returns sane defaults for local development.
func GetDefaultHugotConfig() *HugotConfig {
	model := os.Getenv("MS_EMBEDDING_MODEL")
	if model == "" {
		model = "sentence-transformers/all-MiniLM-L6-v2"
	}
	dim, err := strconv.Atoi(os.Getenv("MS_EMBEDDING_DIM"))
	if err != nil {
		dim = 512
	}

	onnx := os.Getenv("MS_EMBEDDING_MODEL_PATH")
	if onnx == "" {
		onnx = "onnx/model.onnx"
	}
	modelsDir := os.Getenv("MS_EMBEDDING_MODELS_DIR")
	if modelsDir == "" {
		modelsDir = "./models/"
	}

	return &HugotConfig{ModelName: model, Dim: dim, OnnxPath: onnx, ModelsDir: modelsDir}
}

// minimal wrapper to avoid importing pipelines types widely
type featurePipeline struct {
	run func([]string) ([]float32, error)
}

// HugotClient wraps embedders; lazily initializes real Hugot pipeline if available.
type HugotClient struct {
	cfg      *HugotConfig
	session  *hugot.Session
	pipeline *featurePipeline
	initOnce sync.Once
	initErr  error
}

func NewHugotClient(cfg *HugotConfig) *HugotClient { return &HugotClient{cfg: cfg} }

func (c *HugotClient) init() error {
	session, err := hugot.NewGoSession()
	if err != nil {
		return err
	}
	// Download the model.
	downloadOptions := hugot.NewDownloadOptions()
	downloadOptions.OnnxFilePath = c.cfg.OnnxPath
	modelPath, err := hugot.DownloadModel(c.cfg.ModelName, c.cfg.ModelsDir, downloadOptions)
	if err != nil {
		_ = session.Destroy()
		return err
	}
	// Create feature extraction pipeline
	cfg := hugot.FeatureExtractionConfig{ModelPath: modelPath, Name: "music-store-emb"}
	pipeline, err := hugot.NewPipeline(session, cfg)
	if err != nil {
		_ = session.Destroy()
		return err
	}
	c.session = session
	c.pipeline = &featurePipeline{run: func(inputs []string) ([]float32, error) {
		out, err := pipeline.RunPipeline(inputs)
		if err != nil || out == nil || len(out.Embeddings) == 0 {
			return nil, err
		}
		return out.Embeddings[0], nil
	}}
	return nil
}

// Generate returns a DIM-dimensional embedding. Falls back to deterministic if native deps missing.
func (c *HugotClient) Generate(text string) []float32 {
	c.initOnce.Do(func() { c.initErr = c.init() })
	if c.initErr == nil && c.pipeline != nil {
		if v, err := c.pipeline.run([]string{strings.TrimSpace(text)}); err == nil && len(v) > 0 {
			return coerceAndNormalize(v, c.dim())
		}
	}
	return pseudoEmbedding(text, c.dim())
}

func (c *HugotClient) Close() {
	if c.session != nil {
		_ = c.session.Destroy()
	}
}

func (c *HugotClient) dim() int {
	if c == nil || c.cfg == nil || c.cfg.Dim <= 0 {
		return 512
	}
	return c.cfg.Dim
}

func coerceAndNormalize(v []float32, dim int) []float32 {
	out := make([]float32, dim)
	n := dim
	if len(v) < dim {
		n = len(v)
	}
	copy(out[:n], v[:n])
	var sum float64
	for i := 0; i < dim; i++ {
		sum += float64(out[i]) * float64(out[i])
	}
	norm := math.Sqrt(sum)
	if norm == 0 {
		return out
	}
	for i := 0; i < dim; i++ {
		out[i] = float32(float64(out[i]) / norm)
	}
	return out
}

func pseudoEmbedding(text string, dim int) []float32 {
	v := make([]float32, dim)
	for i := 0; i < dim; i++ {
		h := sha256.Sum256([]byte(text + "#" + string(rune(i%997))))
		u := binary.BigEndian.Uint32(h[(i*4)%32 : (i*4)%32+4])
		f := (float64(u%2000000)/1000000.0 - 1.0)
		v[i] = float32(f)
	}
	var sum float64
	for i := 0; i < dim; i++ {
		sum += float64(v[i]) * float64(v[i])
	}
	norm := math.Sqrt(sum)
	if norm == 0 {
		return v
	}
	for i := 0; i < dim; i++ {
		v[i] = float32(float64(v[i]) / norm)
	}
	return v
}
