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
	key := fmt.Sprintf("song:%s", song.Song.Name)
	songJSON, err := json.Marshal(song.Song)
	if err != nil {
		return "Error creating song", err
	}
	if err := r.redisClient.Do(context.Background(), "JSON.SET", key, "$", string(songJSON)).Err(); err != nil {
		return "Error creating song", err
	}
	return "success", nil
}

func (r *songRepository) GetSong(name string) (*model.GetSongResponse, error) {
	key := fmt.Sprintf("song:%s", name)
	songJSON, err := r.redisClient.Do(context.Background(), "JSON.GET", key).Text()
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
	keys, err := r.redisClient.Keys(context.Background(), "song:*").Result()
	if err != nil {
		return nil, err
	}
	var songs []*model.Song
	for _, key := range keys {
		songJSON, err := r.redisClient.Do(context.Background(), "JSON.GET", key).Text()
		if err != nil {
			continue
		}
		var song model.Song
		if err := json.Unmarshal([]byte(songJSON), &song); err != nil {
			continue
		}
		if song.Name == "" {
			continue
		}
		songs = append(songs, &song)
	}
	return &model.GetSongListResponse{Songs: songs}, nil
}

func (r *songRepository) UpdateSong(song *model.UpdateSongRequest) (string, error) {
	name := song.Name
	if name == "" && song.Song.Name != "" {
		name = song.Song.Name
	}
	key := fmt.Sprintf("song:%s", name)
	songJSON, err := json.Marshal(song.Song)
	if err != nil {
		return "Error updating song", err
	}
	if err := r.redisClient.Do(context.Background(), "JSON.SET", key, "$", string(songJSON)).Err(); err != nil {
		return "Error updating song", err
	}
	return "success", nil
}

func (r *songRepository) DeleteSong(name string) (string, error) {
	key := fmt.Sprintf("song:%s", name)
	if err := r.redisClient.Do(context.Background(), "JSON.DEL", key).Err(); err != nil {
		return "Error deleting song", err
	}
	return "success", nil
}
