package handler

import (
	"net/http"

	"wander/backend/internal/service"
)

// Handler holds all HTTP handlers.
type Handler struct {
	UserHandler *UserHandler
}

// NewHandler creates a new Handler with all sub-handlers.
func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(userService),
	}
}

// HealthCheck returns a simple health check response.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"wander-backend"}`))
}
