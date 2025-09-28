package service

import (
	"context"
	"music-store/internal/model"
	"music-store/internal/repository"
	"music-store/utils"
)

type SongService interface {
	CreateSong(ctx context.Context, song *model.CreateSongRequest) (string, error)
	GetSong(ctx context.Context, name string) (*model.GetSongResponse, error)
	GetAllSongs(ctx context.Context) (*model.GetSongListResponse, error)
	UpdateSong(ctx context.Context, song *model.UpdateSongRequest) (string, error)
	DeleteSong(ctx context.Context, name string) (string, error)
	SimilarSongs(ctx context.Context, name string, count int) ([]string, error)
}

type songService struct {
	songRepository repository.SongRepository
	embedder       EmbeddingService
	vectorRepo     repository.VectorRepository
}

func NewSongService(songRepository repository.SongRepository, vectorRepo repository.VectorRepository) SongService {
	// Initialize default embedding service
	embedder := NewEmbeddingService(utils.NewHugotClient(utils.GetDefaultHugotConfig()))
	return &songService{songRepository: songRepository, embedder: embedder, vectorRepo: vectorRepo}
}

func (s *songService) CreateSong(ctx context.Context, song *model.CreateSongRequest) (string, error) {
	song.Song.Embedding = s.embedder.Generate(song.Song.Name)
	if err := s.vectorRepoUpsertSong(song.Song.Name, song.Song.Embedding); err != nil {
		// ignore vector index failure for now
	}
	return s.songRepository.CreateSong(song)
}

func (s *songService) GetSong(ctx context.Context, name string) (*model.GetSongResponse, error) {
	return s.songRepository.GetSong(name)
}

func (s *songService) GetAllSongs(ctx context.Context) (*model.GetSongListResponse, error) {
	return s.songRepository.GetAllSongs()
}

func (s *songService) UpdateSong(ctx context.Context, song *model.UpdateSongRequest) (string, error) {
	song.Song.Embedding = s.embedder.Generate(song.Song.Name)
	if err := s.vectorRepoUpsertSong(song.Song.Name, song.Song.Embedding); err != nil {
		// ignore vector index failure for now
	}
	return s.songRepository.UpdateSong(song)
}

func (s *songService) vectorRepoUpsertSong(name string, emb []float32) error {
	if s.vectorRepo == nil {
		return nil
	}
	return s.vectorRepo.UpsertSong(name, emb)
}

func (s *songService) DeleteSong(ctx context.Context, name string) (string, error) {
	return s.songRepository.DeleteSong(name)
}

func (s *songService) SimilarSongs(ctx context.Context, name string, count int) ([]string, error) {
	songs, err := s.vectorRepo.SimilarSongsByVector(s.embedder.Generate(name), count)
	if err != nil {
		return nil, err
	}
	return songs, nil
}
