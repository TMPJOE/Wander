package main

import (
	"fmt"
	"log"
	"net/http"

	"wander/backend/api"
	"wander/backend/internal/config"
	"wander/backend/internal/handler"
	"wander/backend/internal/middleware"
	"wander/backend/internal/repository"
	"wander/backend/internal/service"
)

func main() {
	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize repositories (in-memory for development).
	userRepo := repository.NewInMemoryUserRepository()

	// Initialize services.
	userService := service.NewUserService(userRepo)

	// Initialize handlers.
	h := handler.NewHandler(userService)

	// Setup routes.
	mux := api.SetupRoutes(h)

	// Apply middleware.
	var server http.Handler
	server = middleware.Recovery(mux)
	server = middleware.Logger(server)
	server = middleware.CORS(cfg.AllowedOrigins)(server)

	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	log.Printf("🚀 Server running on http://%s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
