package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wander/backend/internal/models"
	"wander/backend/internal/service"
	"wander/backend/internal/utils"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes registers user routes on the provided mux.
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", h.List)
	mux.HandleFunc("POST /users", h.Create)
	mux.HandleFunc("GET /users/{id}", h.GetByID)
	mux.HandleFunc("PUT /users/{id}", h.Update)
	mux.HandleFunc("DELETE /users/{id}", h.Delete)
}

// Create handles POST /users.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.service.Create(r.Context(), req)
	if err != nil {
		utils.SendError(w, http.StatusConflict, "Failed to create user", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "User created successfully", user)
}

// GetByID handles GET /users/{id}.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "User retrieved successfully", user)
}

// Update handles PUT /users/{id}.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Failed to update user", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "User updated successfully", user)
}

// Delete handles DELETE /users/{id}.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		utils.SendError(w, http.StatusNotFound, "Failed to delete user", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "User deleted successfully", nil)
}

// List handles GET /users.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context(), 100, 0)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to list users", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Users retrieved successfully", users)
}
