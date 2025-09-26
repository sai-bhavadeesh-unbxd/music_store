package service

import (
	"context"
	"music-store/internal/model"
	"music-store/internal/repository"
)

type SongService interface {
	CreateSong(ctx context.Context, song *model.CreateSongRequest) (string, error)
	GetSong(ctx context.Context, name string) (*model.GetSongResponse, error)
	GetAllSongs(ctx context.Context) (*model.GetSongListResponse, error)
	UpdateSong(ctx context.Context, song *model.UpdateSongRequest) (string, error)
	DeleteSong(ctx context.Context, name string) (string, error)
}

type songService struct {
	songRepository repository.SongRepository
}

func NewSongService(songRepository repository.SongRepository) SongService {
	return &songService{songRepository: songRepository}
}

func (s *songService) CreateSong(ctx context.Context, song *model.CreateSongRequest) (string, error) {
	return s.songRepository.CreateSong(song)
}

func (s *songService) GetSong(ctx context.Context, name string) (*model.GetSongResponse, error) {
	return s.songRepository.GetSong(name)
}

func (s *songService) GetAllSongs(ctx context.Context) (*model.GetSongListResponse, error) {
	return s.songRepository.GetAllSongs()
}

func (s *songService) UpdateSong(ctx context.Context, song *model.UpdateSongRequest) (string, error) {
	return s.songRepository.UpdateSong(song)
}

func (s *songService) DeleteSong(ctx context.Context, name string) (string, error) {
	return s.songRepository.DeleteSong(name)
}
