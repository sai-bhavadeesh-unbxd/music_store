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
	// Marshal the User struct to JSON
	userJSON, err := json.Marshal(user.User)
	if err != nil {
		return "Error marshaling user data", err
	}

	// Store in Redis using namespaced key: user:{id}
	key := fmt.Sprintf("user:%s", user.User.ID)
	_, err = r.redisClient.Set(context.Background(), key, userJSON, 0).Result()
	if err != nil {
		return "Error creating user", err
	}
	return "success", nil
}

func (r *userRepository) GetUser(id string) (*model.GetUserResponse, error) {
	// Use namespaced key: user:{id}
	key := fmt.Sprintf("user:%s", id)
	userJSON, err := r.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON to User struct
	var user model.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, err
	}

	return &model.GetUserResponse{User: &user}, nil
}

func (r *userRepository) GetAllUsers() (*model.GetUserListResponse, error) {
	// Get only user keys using pattern matching: user:*
	keys, err := r.redisClient.Keys(context.Background(), "user:*").Result()
	if err != nil {
		return nil, err
	}

	var users []*model.User
	for _, key := range keys {
		userJSON, err := r.redisClient.Get(context.Background(), key).Result()
		if err != nil {
			continue // Skip keys that can't be retrieved
		}

		// Unmarshal as User struct
		var user model.User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			continue // Skip malformed data
		}

		users = append(users, &user)
	}

	return &model.GetUserListResponse{Users: users}, nil
}

func (r *userRepository) UpdateUser(user *model.UpdateUserRequest) (string, error) {
	// Set the ID from the path parameter if not already set
	if user.User.ID == "" {
		user.User.ID = user.ID
	}

	// Marshal the User struct to JSON
	userJSON, err := json.Marshal(user.User)
	if err != nil {
		return "Error marshaling user data", err
	}

	// Use namespaced key: user:{id}
	key := fmt.Sprintf("user:%s", user.ID)
	_, err = r.redisClient.Set(context.Background(), key, userJSON, 0).Result()
	if err != nil {
		return "Error updating user", err
	}
	return "success", nil
}

func (r *userRepository) DeleteUser(id string) (string, error) {
	// Use namespaced key: user:{id}
	key := fmt.Sprintf("user:%s", id)
	_, err := r.redisClient.Del(context.Background(), key).Result()
	if err != nil {
		return "Error deleting user", err
	}
	return "success", nil
}
