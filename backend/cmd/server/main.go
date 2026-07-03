package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"wander/backend/api"
	"wander/backend/internal/config"
	"wander/backend/internal/handler"
	"wander/backend/internal/middleware"
	"wander/backend/internal/repository"
	"wander/backend/internal/service"
)

func main() {
	// Initialize structured logging (slog).
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	slog.Info("starting Wander API server...")

	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize DB Connection pool.
	dbPool, err := config.ConnectDB(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Run DB migrations.
	cwd, _ := os.Getwd()
	migrationsDir := filepath.Join(cwd, "backend", "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = filepath.Join(cwd, "migrations") // fallback if already in backend dir
	}
	slog.Info("running database migrations", "path", migrationsDir)
	if err := config.RunMigrations(dbPool, migrationsDir); err != nil {
		slog.Error("failed to run database migrations", "error", err)
		os.Exit(1)
	}

	// Initialize pgx repositories.
	userRepo := repository.NewPgUserRepository(dbPool)
	categoryRepo := repository.NewPgCategoryRepository(dbPool)
	tourRepo := repository.NewPgTourRepository(dbPool)
	scheduleRepo := repository.NewPgTourScheduleRepository(dbPool)
	bookingRepo := repository.NewPgBookingRepository(dbPool)
	reviewRepo := repository.NewPgReviewRepository(dbPool)
	favoriteRepo := repository.NewPgFavoriteRepository(dbPool)
	messageRepo := repository.NewPgMessageRepository(dbPool)

	// Initialize services.
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tourService := service.NewTourService(tourRepo)
	scheduleService := service.NewTourScheduleService(scheduleRepo, tourRepo)
	bookingService := service.NewBookingService(bookingRepo, scheduleRepo, tourRepo)
	reviewService := service.NewReviewService(reviewRepo, bookingRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo)
	messageService := service.NewMessageService(messageRepo)

	// Initialize handlers.
	h := handler.NewHandler(
		authService,
		userService,
		categoryService,
		tourService,
		scheduleService,
		bookingService,
		reviewService,
		favoriteService,
		messageService,
	)

	// Setup routes.
	mux := api.SetupRoutes(h, cfg.JWTSecret)

	// Apply middleware.
	var server http.Handler
	server = middleware.Recovery(mux)
	server = middleware.Logger(server)
	server = middleware.CORS(cfg.AllowedOrigins)(server)

	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	slog.Info("🚀 Server running", "addr", fmt.Sprintf("http://%s", addr))
	if err := http.ListenAndServe(addr, server); err != nil {
		slog.Error("server execution failed", "error", err)
		os.Exit(1)
	}
}
