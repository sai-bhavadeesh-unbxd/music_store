package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Address  string
	Password string
	Database int
}

// GetDefaultRedisConfig returns default Redis configuration
func GetDefaultRedisConfig() *RedisConfig {
	address := os.Getenv("REDIS_ADDR")
	if address == "" {
		address = "localhost:6379" // Default for local development
	}

	password := os.Getenv("REDIS_PASSWORD")
	return &RedisConfig{
		Address:  address,
		Password: password,
		Database: 0, // Default database
	}
}

// InitRedis initializes Redis connection
func InitRedis(config *RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.Address,
		Password:     config.Password,
		DB:           config.Database,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	// Test the connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("failed to connect to Redis: %v", err)
		return err
	}

	log.Println("Successfully connected to Redis")
	return nil
}

// InitRedisWithDefaults initializes Redis with default configuration
func InitRedisWithDefaults() error {
	config := GetDefaultRedisConfig()
	return InitRedis(config)
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	if RedisClient == nil {
		log.Println("Redis client is not initialized. Call InitRedis first.")
	}
	return RedisClient
}

// CloseRedis closes Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Printf("failed to close Redis connection: %v", err)
			return err
		}
		log.Println("Redis connection closed")
	}
	return nil
}

// GetContext returns the context for Redis operations
func GetContext() context.Context {
	return ctx
}
