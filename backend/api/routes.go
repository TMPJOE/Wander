package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"wander/backend/internal/handler"
	"wander/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes configures all application routes using chi.
func SetupRoutes(h *handler.Handler, jwtSecret string) *chi.Mux {
	r := chi.NewRouter()

	// Global chi middleware: request id, recover from panics, normalize slashes,
	// and clean trailing/double slashes.
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(chimw.StripSlashes)

	// Public routes
	r.Get("/api/v1/health", h.HealthCheck)
	r.Post("/api/v1/auth/register", h.AuthHandler.Register)
	r.Post("/api/v1/auth/login", h.AuthHandler.Login)

	// Public GET routes for tours, categories, schedules, reviews, and users
	r.Get("/api/v1/categories", h.CategoryHandler.List)
	r.Get("/api/v1/tours", h.TourHandler.List)
	r.Get("/api/v1/tours/{id}", h.TourHandler.GetByID)
	r.Get("/api/v1/tours/{tourId}/schedules", h.ScheduleHandler.ListByTourID)
	r.Get("/api/v1/tours/{tourId}/reviews", h.ReviewHandler.ListByTour)
	r.Get("/api/v1/users/{id}", h.UserHandler.GetByID)

	// Auth middleware reused across protected groups.
	authMiddleware := middleware.Auth(jwtSecret)
	guideOnly := middleware.RequireRole("guide")

	// Upload route — any authenticated user (guide creates tours; travelers
	// may upload avatars later). Mounted at /api/v1/uploads.
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/api/v1/uploads", h.UploadHandler.UploadImage)
	})

	// Traveler / shared authenticated routes.
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		// Profile
		r.Get("/api/v1/users/me", h.UserHandler.GetMe)
		r.Put("/api/v1/users/me", h.UserHandler.UpdateMe)

		// Bookings (shared path; mutating guide-only actions are registered
		// below on the same route via a separate guide group).
		r.Post("/api/v1/bookings", h.BookingHandler.Create)
		r.Get("/api/v1/bookings", h.BookingHandler.List)
		r.Get("/api/v1/bookings/{id}", h.BookingHandler.GetByID)
		r.Patch("/api/v1/bookings/{id}/cancel", h.BookingHandler.Cancel)

		// Payments
		r.Post("/api/v1/payments/bookings/{id}/intent", h.PaymentHandler.CreateIntent)
		r.Post("/api/v1/payments/bookings/{id}/confirm", h.PaymentHandler.Confirm)

		// Reviews (any authenticated user may review)
		r.Post("/api/v1/tours/{tourId}/reviews", h.ReviewHandler.Create)
		r.Put("/api/v1/reviews/{id}", h.ReviewHandler.Update)
		r.Get("/api/v1/reviews/me", h.ReviewHandler.ListByUser)

		// Favorites
		r.Get("/api/v1/favorites", h.FavoriteHandler.List)
		r.Post("/api/v1/favorites/{tourId}", h.FavoriteHandler.Add)
		r.Delete("/api/v1/favorites/{tourId}", h.FavoriteHandler.Remove)

		// Messages
		r.Get("/api/v1/messages/conversations", h.MessageHandler.ListConversations)
		r.Get("/api/v1/messages/stream", h.MessageHandler.StreamMessages)
		r.Get("/api/v1/messages/{userId}", h.MessageHandler.ListMessages)
		r.Post("/api/v1/messages/{userId}", h.MessageHandler.Create)
	})

	// Guide-only authenticated routes.
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(guideOnly)

		// Tours (mutations)
		r.Post("/api/v1/tours", h.TourHandler.Create)
		r.Put("/api/v1/tours/{id}", h.TourHandler.Update)
		r.Delete("/api/v1/tours/{id}", h.TourHandler.Delete)
		r.Get("/api/v1/guide/tours", h.TourHandler.ListMyTours)
		r.Get("/api/v1/guide/stats", h.TourHandler.GetStats)

		// Schedules
		r.Post("/api/v1/schedules", h.ScheduleHandler.Create)
		r.Put("/api/v1/schedules/{id}", h.ScheduleHandler.Update)
		r.Delete("/api/v1/schedules/{id}", h.ScheduleHandler.Delete)

		// Booking confirmations / lifecycle
		r.Get("/api/v1/guide/bookings", h.BookingHandler.GetGuideBookings)
		r.Get("/api/v1/guide/earnings", h.BookingHandler.GetEarnings)
		r.Patch("/api/v1/bookings/{id}/confirm", h.BookingHandler.Confirm)
		r.Patch("/api/v1/bookings/{id}/complete", h.BookingHandler.Complete)
		r.Patch("/api/v1/bookings/{id}/reject", h.BookingHandler.Reject)
	})

	return r
}

// SetupStaticRoutes registers uploaded-image serving and the SPA static-file
// fallback onto the same chi router used by the API endpoints.
func SetupStaticRoutes(r chi.Router, uploadsDir, distDir string) {
	// Serve uploaded images at /uploads/.
	if err := os.MkdirAll(uploadsDir, 0o755); err == nil {
		r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))
	}

	// Serve frontend static files (production build) with SPA fallback.
	if _, err := os.Stat(distDir); err == nil {
		fs := http.FileServer(http.Dir(distDir))

		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			filePath := filepath.Join(distDir, r.URL.Path)
			if _, err := os.Stat(filePath); err == nil {
				fs.ServeHTTP(w, r)
				return
			}

			// Fallback: serve index.html for SPA client-side routing.
			indexPath := filepath.Join(distDir, "index.html")
			f, err := os.Open(indexPath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			defer f.Close()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = io.Copy(w, f)
		})
	}
}
