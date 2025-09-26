package service

import (
	"context"
	"music-store/internal/model"
	"music-store/internal/repository"
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
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) CreateUser(ctx context.Context, user *model.CreateUserRequest) (string, error) {
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

	// Append and persist
	user.LikedSongs = append(user.LikedSongs, songName)

	updateReq := &model.UpdateUserRequest{ID: userID, User: *user}
	if _, err := s.userRepository.UpdateUser(updateReq); err != nil {
		return "Error updating user", err
	}
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

	updateReq := &model.UpdateUserRequest{ID: userID, User: *user}
	if _, err := s.userRepository.UpdateUser(updateReq); err != nil {
		return "Error updating user", err
	}
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
