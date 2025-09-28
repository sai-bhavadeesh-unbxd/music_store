package service

import (
	"context"
	"music-store/internal/model"
	"music-store/internal/repository"
	"music-store/utils"
	"strings"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.CreateUserRequest) (string, error)
	GetUser(ctx context.Context, id string) (*model.GetUserResponse, error)
	GetAllUsers(ctx context.Context) (*model.GetUserListResponse, error)
	UpdateUser(ctx context.Context, user *model.UpdateUserRequest) (string, error)
	DeleteUser(ctx context.Context, id string) (string, error)
	LikeSong(ctx context.Context, userID, songName string) (string, error)
	UnlikeSong(ctx context.Context, userID, songName string) (string, error)
	GetLikedSongs(ctx context.Context, userID string) ([]string, error)
	GetRecommendedSongs(ctx context.Context, userID string, count int) ([]string, error)
}

type userService struct {
	userRepository repository.UserRepository
	embedder       EmbeddingService
	vectorRepo     repository.VectorRepository
}

func NewUserService(userRepository repository.UserRepository, vectorRepo repository.VectorRepository) UserService {
	embedder := NewEmbeddingService(utils.NewHugotClient(utils.GetDefaultHugotConfig()))
	return &userService{userRepository: userRepository, embedder: embedder, vectorRepo: vectorRepo}
}

func (s *userService) CreateUser(ctx context.Context, user *model.CreateUserRequest) (string, error) {
	// If user embedding is absent, initialize (e.g., zero vector)
	if len(user.User.Embedding) == 0 {
		user.User.Embedding = make([]float32, s.embedder.Dim())
	}
	if err := s.vectorRepoUpsertUser(user.User.ID, user.User.Embedding); err != nil {
		// ignore vector index failure for now
	}
	return s.userRepository.CreateUser(user)
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.GetUserResponse, error) {
	return s.userRepository.GetUser(id)
}

func (s *userService) GetAllUsers(ctx context.Context) (*model.GetUserListResponse, error) {
	return s.userRepository.GetAllUsers()
}

func (s *userService) UpdateUser(ctx context.Context, user *model.UpdateUserRequest) (string, error) {
	return s.userRepository.UpdateUser(user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) (string, error) {
	return s.userRepository.DeleteUser(id)
}

func (s *userService) LikeSong(ctx context.Context, userID, songName string) (string, error) {
	// Fetch user
	userResp, err := s.userRepository.GetUser(userID)
	if err != nil {
		return "Error getting user", err
	}
	if userResp == nil || userResp.User == nil {
		return "Error getting user", nil
	}

	user := userResp.User

	// Initialize LikedSongs if nil
	if user.LikedSongs == nil {
		user.LikedSongs = []string{}
	}

	// Check duplicate
	for _, liked := range user.LikedSongs {
		if liked == songName {
			return "Song already liked", nil
		}
	}

	// Append and recompute embedding from all liked songs (actual song embeddings)
	user.LikedSongs = append(user.LikedSongs, songName)
	if vec, err := s.computeEmbeddingFromLikes(user.LikedSongs); err == nil {
		user.Embedding = vec
	} else {
		user.Embedding = make([]float32, s.embedder.Dim())
	}

	updateReq := &model.UpdateUserRequest{ID: userID, User: *user}
	if _, err := s.userRepository.UpdateUser(updateReq); err != nil {
		return "Error updating user", err
	}
	_ = s.vectorRepoUpsertUser(userID, user.Embedding)
	return "success", nil
}

func (s *userService) UnlikeSong(ctx context.Context, userID, songName string) (string, error) {
	// Fetch user
	userResp, err := s.userRepository.GetUser(userID)
	if err != nil {
		return "Error getting user", err
	}
	if userResp == nil || userResp.User == nil {
		return "Error getting user", nil
	}

	user := userResp.User

	// Initialize LikedSongs if nil
	if user.LikedSongs == nil {
		user.LikedSongs = []string{}
	}

	// Remove if present
	updated := make([]string, 0, len(user.LikedSongs))
	found := false
	for _, liked := range user.LikedSongs {
		if liked != songName {
			updated = append(updated, liked)
		} else {
			found = true
		}
	}

	if !found {
		return "Song was not liked", nil
	}

	user.LikedSongs = updated

	// Recompute embedding from remaining likes using actual song embeddings
	if vec, err := s.computeEmbeddingFromLikes(user.LikedSongs); err == nil {
		user.Embedding = vec
	} else {
		user.Embedding = make([]float32, s.embedder.Dim())
	}

	updateReq := &model.UpdateUserRequest{ID: userID, User: *user}
	if _, err := s.userRepository.UpdateUser(updateReq); err != nil {
		return "Error updating user", err
	}
	_ = s.vectorRepoUpsertUser(userID, user.Embedding)
	return "success", nil
}

func (s *userService) GetLikedSongs(ctx context.Context, userID string) ([]string, error) {
	userResp, err := s.userRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return []string{}, nil
	}
	if userResp.User.LikedSongs == nil {
		return []string{}, nil
	}
	return userResp.User.LikedSongs, nil
}

func (s *userService) GetRecommendedSongs(ctx context.Context, userID string, count int) ([]string, error) {
	userResp, err := s.userRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return []string{}, nil
	}
	songs, err := s.vectorRepo.SimilarSongsByVector(userResp.User.Embedding, count)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (s *userService) vectorRepoUpsertUser(id string, emb []float32) error {
	if s.vectorRepo == nil {
		return nil
	}
	return s.vectorRepo.UpsertUser(id, emb)
}

// computeEmbeddingFromLikes averages embeddings of all liked songs.
func (s *userService) computeEmbeddingFromLikes(liked []string) ([]float32, error) {
	dim := s.embedder.Dim()
	if len(liked) == 0 {
		return make([]float32, dim), nil
	}
	// Join liked song names into a single string and embed once
	text := strings.Join(liked, " ")
	vec := s.embedder.Generate(text)
	// Coerce to expected dim if needed
	if len(vec) != dim {
		out := make([]float32, dim)
		n := dim
		if len(vec) < dim {
			n = len(vec)
		}
		copy(out[:n], vec[:n])
		vec = out
	}
	return vec, nil
}
