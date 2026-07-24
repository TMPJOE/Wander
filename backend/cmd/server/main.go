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
	"wander/backend/migrations"
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

	// Run DB migrations (embedded at compile time — no cwd/path guessing).
	cwd, _ := os.Getwd()
	slog.Info("running database migrations")
	if err := config.RunMigrations(dbPool, migrations.FS); err != nil {
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
	paymentService := service.NewPaymentService(bookingRepo, cfg.StripeSecretKey, cfg.StripePublishableKey)
	bookingService := service.NewBookingService(bookingRepo, scheduleRepo, tourRepo, paymentService)
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
		paymentService,
	)

	// Wire upload handler (saved under backend/uploads, served at /uploads).
	var uploadsDir string
	if _, err := os.Stat(filepath.Join(cwd, "backend")); err == nil {
		uploadsDir = filepath.Join(cwd, "backend", "uploads")
	} else {
		uploadsDir = filepath.Join(cwd, "uploads") // fallback if running from backend/
	}
	h.UploadHandler = handler.NewUploadHandler(uploadsDir)

	// Setup routes (chi router returned as http.Handler).
	r := api.SetupRoutes(h, cfg.JWTSecret)

	// Serve uploaded images at /uploads/.
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		slog.Error("failed to create uploads dir", "error", err)
	}

	// Serve frontend static files (production build).
	// Look for dist/ relative to cwd (project root or backend dir).
	distDir := filepath.Join(cwd, "frontend", "dist")
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		distDir = filepath.Join(cwd, "..", "frontend", "dist") // fallback if running from backend/
	}

	if info, err := os.Stat(distDir); err == nil && info.IsDir() {
		slog.Info("serving frontend static files", "path", distDir)
	} else {
		slog.Warn("frontend dist/ not found, serving API only", "checked", distDir)
		distDir = "" // signals SetupStaticRoutes to skip the SPA handler
	}

	// Register static (uploads + SPA fallback) routes onto the same chi router.
	api.SetupStaticRoutes(r, uploadsDir, distDir)

	// Apply middleware.
	var server http.Handler
	server = middleware.Recovery(r)
	server = middleware.Logger(server)
	server = middleware.CORS(cfg.AllowedOrigins)(server)

	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	slog.Info("🚀 Server running", "addr", fmt.Sprintf("http://%s", addr))
	if err := http.ListenAndServe(addr, server); err != nil {
		slog.Error("server execution failed", "error", err)
		os.Exit(1)
	}
}
