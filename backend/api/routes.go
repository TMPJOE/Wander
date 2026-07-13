package api

import (
	"net/http"
	"strings"

	"wander/backend/internal/handler"
	"wander/backend/internal/middleware"
)

// SetupRoutes configures all application routes, wrapping with Auth where required.
func SetupRoutes(h *handler.Handler, jwtSecret string) *http.ServeMux {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /api/v1/health", h.HealthCheck)
	mux.HandleFunc("POST /api/v1/auth/register", h.AuthHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", h.AuthHandler.Login)

	mux.HandleFunc("GET /api/v1/categories", h.CategoryHandler.List)
	mux.HandleFunc("GET /api/v1/tours", h.TourHandler.List)
	mux.HandleFunc("GET /api/v1/tours/{id}", h.TourHandler.GetByID)
	mux.HandleFunc("GET /api/v1/tours/{tourId}/schedules", h.ScheduleHandler.ListByTourID)
	mux.HandleFunc("GET /api/v1/tours/{tourId}/reviews", h.ReviewHandler.ListByTour)
	mux.HandleFunc("GET /api/v1/users/{id}", h.UserHandler.GetByID)

	// Auth middleware
	authMiddleware := middleware.Auth(jwtSecret)

	// Upload route — any authenticated user (guide creates tours; travelers
	// may upload avatars later). Mounted at /api/v1/uploads.
	uploadMux := http.NewServeMux()
	uploadMux.HandleFunc("POST /uploads", h.UploadHandler.UploadImage)
	uploadGroup := authMiddleware(http.StripPrefix("/api/v1", uploadMux))
	mux.Handle("/api/v1/uploads", uploadGroup)
	mux.Handle("/api/v1/uploads/", uploadGroup)

	// Private Traveler/Shared Routes
	travelerMux := http.NewServeMux()
	travelerMux.HandleFunc("GET /users/me", h.UserHandler.GetMe)
	travelerMux.HandleFunc("PUT /users/me", h.UserHandler.UpdateMe)

	travelerMux.HandleFunc("POST /bookings", h.BookingHandler.Create)
	travelerMux.HandleFunc("GET /bookings", h.BookingHandler.List)
	travelerMux.HandleFunc("GET /bookings/{id}", h.BookingHandler.GetByID)
	travelerMux.HandleFunc("PATCH /bookings/{id}/cancel", h.BookingHandler.Cancel)

	travelerMux.HandleFunc("POST /payments/bookings/{id}/intent", h.PaymentHandler.CreateIntent)
	travelerMux.HandleFunc("POST /payments/bookings/{id}/confirm", h.PaymentHandler.Confirm)

	travelerMux.HandleFunc("POST /tours/{tourId}/reviews", h.ReviewHandler.Create)
	travelerMux.HandleFunc("PUT /reviews/{id}", h.ReviewHandler.Update)
	travelerMux.HandleFunc("GET /reviews/me", h.ReviewHandler.ListByUser)

	travelerMux.HandleFunc("GET /favorites", h.FavoriteHandler.List)
	travelerMux.HandleFunc("POST /favorites/{tourId}", h.FavoriteHandler.Add)
	travelerMux.HandleFunc("DELETE /favorites/{tourId}", h.FavoriteHandler.Remove)

	travelerMux.HandleFunc("GET /messages/conversations", h.MessageHandler.ListConversations)
	travelerMux.HandleFunc("GET /messages/{userId}", h.MessageHandler.ListMessages)
	travelerMux.HandleFunc("POST /messages/{userId}", h.MessageHandler.Create)
	travelerMux.HandleFunc("GET /messages/stream", h.MessageHandler.StreamMessages)

	// Guide-only Routes
	guideMux := http.NewServeMux()
	guideMux.HandleFunc("POST /tours", h.TourHandler.Create)
	guideMux.HandleFunc("PUT /tours/{id}", h.TourHandler.Update)
	guideMux.HandleFunc("DELETE /tours/{id}", h.TourHandler.Delete)
	guideMux.HandleFunc("GET /guide/tours", h.TourHandler.ListMyTours)
	guideMux.HandleFunc("GET /guide/stats", h.TourHandler.GetStats)

	guideMux.HandleFunc("POST /schedules", h.ScheduleHandler.Create)
	guideMux.HandleFunc("PUT /schedules/{id}", h.ScheduleHandler.Update)
	guideMux.HandleFunc("DELETE /schedules/{id}", h.ScheduleHandler.Delete)

	guideMux.HandleFunc("PATCH /bookings/{id}/confirm", h.BookingHandler.Confirm)
	guideMux.HandleFunc("PATCH /bookings/{id}/complete", h.BookingHandler.Complete)
	guideMux.HandleFunc("PATCH /bookings/{id}/reject", h.BookingHandler.Reject)

	// Guide endpoint routing group (declared early so guide-only PATCH routes on
	// shared paths like /bookings/ can dispatch into it below).
	guideGroup := authMiddleware(middleware.RequireRole("guide")(http.StripPrefix("/api/v1", guideMux)))

	// Register sub-muxes
	travelerGroup := authMiddleware(http.StripPrefix("/api/v1", travelerMux))
	mux.Handle("/api/v1/users/", travelerGroup)
	mux.Handle("/api/v1/users", travelerGroup)
	mux.Handle("/api/v1/bookings/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// PATCH confirm/complete/reject are guide-only; everything else is traveler.
		if r.Method == http.MethodPatch && (strings.HasSuffix(r.URL.Path, "/confirm") || strings.HasSuffix(r.URL.Path, "/complete") || strings.HasSuffix(r.URL.Path, "/reject")) {
			guideGroup.ServeHTTP(w, r)
			return
		}
		travelerGroup.ServeHTTP(w, r)
	}))
	mux.Handle("/api/v1/bookings", travelerGroup)
	mux.Handle("/api/v1/favorites/", travelerGroup)
	mux.Handle("/api/v1/favorites", travelerGroup)
	mux.Handle("/api/v1/messages/", travelerGroup)
	mux.Handle("/api/v1/messages", travelerGroup)
	mux.Handle("/api/v1/reviews/", travelerGroup)
	mux.Handle("/api/v1/reviews", travelerGroup)
	mux.Handle("/api/v1/payments/", travelerGroup)
	mux.Handle("/api/v1/payments", travelerGroup)

	mux.Handle("/api/v1/tours/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Public GET falls back to root routing
			http.DefaultServeMux.ServeHTTP(w, r)
			return
		}
		// POST reviews should go to travelerGroup (any authenticated user can review)
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/reviews") {
			travelerGroup.ServeHTTP(w, r)
			return
		}
		// All other POST/PATCH/DELETE on tours are guide-only
		guideGroup.ServeHTTP(w, r)
	}))
	mux.Handle("/api/v1/tours", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.DefaultServeMux.ServeHTTP(w, r)
			return
		}
		guideGroup.ServeHTTP(w, r)
	}))
	mux.Handle("/api/v1/schedules/", guideGroup)
	mux.Handle("/api/v1/schedules", guideGroup)
	mux.Handle("/api/v1/guide/", guideGroup)
	mux.Handle("/api/v1/guide", guideGroup)

	return mux
}
