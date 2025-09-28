package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"music-store/internal/model"

	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	CreateUser(user *model.CreateUserRequest) (string, error)
	GetUser(id string) (*model.GetUserResponse, error)
	GetAllUsers() (*model.GetUserListResponse, error)
	UpdateUser(user *model.UpdateUserRequest) (string, error)
	DeleteUser(id string) (string, error)
}

type userRepository struct {
	redisClient *redis.Client
}

func NewUserRepository(redisClient *redis.Client) UserRepository {
	return &userRepository{redisClient: redisClient}
}

func (r *userRepository) CreateUser(user *model.CreateUserRequest) (string, error) {
	key := fmt.Sprintf("user:%s", user.User.ID)
	userJSON, err := json.Marshal(user.User)
	if err != nil {
		return "Error creating user", err
	}
	if err := r.redisClient.Do(context.Background(), "JSON.SET", key, "$", string(userJSON)).Err(); err != nil {
		return "Error creating user", err
	}
	return "success", nil
}

func (r *userRepository) GetUser(id string) (*model.GetUserResponse, error) {
	key := fmt.Sprintf("user:%s", id)
	// Get entire JSON document (omit path to avoid $ array wrapper)
	userJSON, err := r.redisClient.Do(context.Background(), "JSON.GET", key).Text()
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, err
	}
	return &model.GetUserResponse{User: &user}, nil
}

func (r *userRepository) GetAllUsers() (*model.GetUserListResponse, error) {
	keys, err := r.redisClient.Keys(context.Background(), "user:*").Result()
	if err != nil {
		return nil, err
	}
	var users []*model.User
	for _, key := range keys {
		userJSON, err := r.redisClient.Do(context.Background(), "JSON.GET", key).Text()
		if err != nil {
			continue
		}
		var user model.User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			continue
		}
		users = append(users, &user)
	}
	return &model.GetUserListResponse{Users: users}, nil
}

func (r *userRepository) UpdateUser(user *model.UpdateUserRequest) (string, error) {
	if user.User.ID == "" {
		user.User.ID = user.ID
	}
	key := fmt.Sprintf("user:%s", user.ID)
	userJSON, err := json.Marshal(user.User)
	if err != nil {
		return "Error updating user", err
	}
	if err := r.redisClient.Do(context.Background(), "JSON.SET", key, "$", string(userJSON)).Err(); err != nil {
		return "Error updating user", err
	}
	return "success", nil
}

func (r *userRepository) DeleteUser(id string) (string, error) {
	key := fmt.Sprintf("user:%s", id)
	if err := r.redisClient.Do(context.Background(), "JSON.DEL", key).Err(); err != nil {
		return "Error deleting user", err
	}
	return "success", nil
}
