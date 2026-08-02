package main

import (
	"context"
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
	"wander/backend/internal/storage"
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

	// Wire storage provider.
	//
	// Two storage backends are selectable via STORAGE_DRIVER:
	//   - "local" (default): files land on disk and are served at /uploads/*.
	//   - "s3": files stream into an S3-compatible bucket; returned URLs are
	//     absolute, so the Go server no longer serves /uploads/* in this mode.
	// The provider is created early because it is shared between the upload
	// handler (persisting new uploads) and the tour service (cleaning up
	// orphaned files when tours/images are deleted).
	var (
		uploadsDir string
		provider   storage.Provider
	)

	switch cfg.Storage.Driver {
	case "s3":
		p, err := storage.NewS3Provider(context.Background(), storage.S3Options{
			Bucket:         cfg.Storage.S3Bucket,
			Region:         cfg.Storage.S3Region,
			Endpoint:       cfg.Storage.S3Endpoint,
			AccessKey:      cfg.Storage.S3AccessKey,
			SecretKey:      cfg.Storage.S3SecretKey,
			ForcePathStyle: cfg.Storage.S3ForcePathStyle,
			PublicBaseURL:  cfg.Storage.S3PublicBaseURL,
		})
		if err != nil {
			slog.Error("failed to init S3 storage provider", "error", err)
			os.Exit(1)
		}
		provider = p
		slog.Info("storage provider", "driver", "s3", "bucket", cfg.Storage.S3Bucket)
		// In S3 mode we do NOT register the /uploads/* static route below.

	case "local", "":
		locallyRelative := cfg.Storage.UploadsDir
		if _, err := os.Stat(filepath.Join(cwd, "backend")); err == nil {
			uploadsDir = filepath.Join(cwd, "backend", locallyRelative)
		} else {
			uploadsDir = filepath.Join(cwd, locallyRelative) // fallback if running from backend/
		}
		p, err := storage.NewLocalProvider(uploadsDir, cfg.Storage.PublicBaseURL)
		if err != nil {
			slog.Error("failed to init local storage provider", "error", err)
			os.Exit(1)
		}
		provider = p
		slog.Info("storage provider", "driver", "local", "dir", uploadsDir, "base", cfg.Storage.PublicBaseURL)

	default:
		slog.Error("unknown STORAGE_DRIVER", "driver", cfg.Storage.Driver)
		os.Exit(1)
	}

	// Initialize services.
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tourService := service.NewTourService(tourRepo, provider)
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

	// Wire upload handler with the shared provider.
	h.UploadHandler = handler.NewUploadHandler(provider)

	// Setup routes (chi router returned as http.Handler).
	r := api.SetupRoutes(h, cfg.JWTSecret)

	// In local storage mode ensure the uploads directory exists and serve
	// it at /uploads/*. In S3 mode uploadsDir is empty and we skip this so
	// bucket URLs are served by the bucket/CDN, not the Go server.
	if uploadsDir != "" {
		if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
			slog.Error("failed to create uploads dir", "error", err)
		}
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

	// Register static routes. Pass an empty uploadsDir in S3 mode so
	// SetupStaticRoutes skips the /uploads/* FileServer mount.
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
