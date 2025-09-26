package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"music-store/internal/model"

	"github.com/redis/go-redis/v9"
)

type SongRepository interface {
	CreateSong(song *model.CreateSongRequest) (string, error)
	GetSong(name string) (*model.GetSongResponse, error)
	GetAllSongs() (*model.GetSongListResponse, error)
	UpdateSong(song *model.UpdateSongRequest) (string, error)
	DeleteSong(name string) (string, error)
}

type songRepository struct {
	redisClient *redis.Client
}

func NewSongRepository(redisClient *redis.Client) SongRepository {
	return &songRepository{redisClient: redisClient}
}

func (r *songRepository) CreateSong(song *model.CreateSongRequest) (string, error) {
	// Marshal the Song struct to JSON
	songJSON, err := json.Marshal(&song.Song)
	if err != nil {
		return "Error marshaling song data", err
	}

	// Store in Redis using namespaced key: song:{name}
	key := fmt.Sprintf("song:%s", song.Song.Name)
	_, err = r.redisClient.Set(context.Background(), key, songJSON, 0).Result()
	if err != nil {
		return "Error creating song", err
	}
	return "success", nil
}

func (r *songRepository) GetSong(name string) (*model.GetSongResponse, error) {
	// Use namespaced key: song:{name}
	key := fmt.Sprintf("song:%s", name)
	songJSON, err := r.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var song model.Song
	if err := json.Unmarshal([]byte(songJSON), &song); err != nil {
		return nil, err
	}
	return &model.GetSongResponse{Song: &song}, nil
}

func (r *songRepository) GetAllSongs() (*model.GetSongListResponse, error) {
	// Get only song keys using pattern matching: song:*
	keys, err := r.redisClient.Keys(context.Background(), "song:*").Result()
	if err != nil {
		return nil, err
	}

	var songs []*model.Song
	for _, key := range keys {
		songJSON, err := r.redisClient.Get(context.Background(), key).Result()
		if err != nil {
			continue // Skip keys that can't be retrieved
		}

		// Unmarshal as Song struct
		var song model.Song
		if err := json.Unmarshal([]byte(songJSON), &song); err != nil {
			continue // Skip malformed data
		}

		// Validate that this is actually song data (songs must have a name)
		if song.Name == "" {
			continue // Skip data without song name (probably user data)
		}

		songs = append(songs, &song)
	}

	return &model.GetSongListResponse{Songs: songs}, nil
}

func (r *songRepository) UpdateSong(song *model.UpdateSongRequest) (string, error) {
	// Use the Name from the path parameter if the song name is empty
	name := song.Name
	if name == "" && song.Song.Name != "" {
		name = song.Song.Name
	}

	// Marshal the Song struct to JSON
	songJSON, err := json.Marshal(&song.Song)
	if err != nil {
		return "Error marshaling song data", err
	}

	// Store in Redis using namespaced key: song:{name}
	key := fmt.Sprintf("song:%s", name)
	_, err = r.redisClient.Set(context.Background(), key, songJSON, 0).Result()
	if err != nil {
		return "Error updating song", err
	}
	return "success", nil
}

func (r *songRepository) DeleteSong(name string) (string, error) {
	// Use namespaced key: song:{name}
	key := fmt.Sprintf("song:%s", name)
	_, err := r.redisClient.Del(context.Background(), key).Result()
	if err != nil {
		return "Error deleting song", err
	}
	return "success", nil
}
