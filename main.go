package main

import (
	"log"
	"music-store/internal/controller"
	"music-store/internal/repository"
	"music-store/internal/service"
	"music-store/utils"

	"github.com/unbxd/go-base/kit/transport/http"
)

func main() {
	// Initialize Redis connection with default configuration
	err := utils.InitRedisWithDefaults()
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Get Redis client for use in your application
	redisClient := utils.GetRedisClient()
	if redisClient == nil {
		log.Fatal("Redis client is nil")
	}

	// Initialize dependencies
	// Repositories
	userRepo := repository.NewUserRepository(redisClient)
	songRepo := repository.NewSongRepository(redisClient)
	vecRepo := repository.NewVectorRepository(redisClient)

	// Services (inject songRepo into userSvc for embedding recompute)
	userSvc := service.NewUserService(userRepo, vecRepo)
	songSvc := service.NewSongService(songRepo, vecRepo)

	// Controllers
	userController := controller.NewUserController(userSvc)
	songController := controller.NewSongController(songSvc)

	// Initialize HTTP transport
	transport, err := http.NewTransport("0.0.0.0", "8080")
	if err != nil {
		log.Fatalf("Failed to initialize HTTP transport: %v", err)
	}

	// Bind user routes
	userController.Bind(transport, []http.HandlerOption{})
	songController.Bind(transport, []http.HandlerOption{})

	// Close Redis connection when done
	defer func() {
		if err := utils.CloseRedis(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}()

	log.Println("Music Store application started successfully!")

	// Start the HTTP server
	log.Println("Starting HTTP server on :8080...")
	if err := transport.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
