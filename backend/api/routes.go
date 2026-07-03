package api

import (
	"net/http"

	"wander/backend/internal/handler"
)

// SetupRoutes configures all application routes.
func SetupRoutes(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", h.HealthCheck)

	// User routes
	h.UserHandler.RegisterRoutes(mux)

	return mux
}
