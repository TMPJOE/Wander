package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wander/backend/internal/middleware"
	"wander/backend/internal/models"
	"wander/backend/internal/service"
	"wander/backend/internal/utils"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	// Group under prefix handled externally or register fully
	mux.HandleFunc("GET /users/me", h.GetMe)
	mux.HandleFunc("PUT /users/me", h.UpdateMe)
	mux.HandleFunc("GET /users/{id}", h.GetByID)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Usuario no encontrado", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Perfil recuperado", user)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	user, err := h.service.Update(r.Context(), userID, req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al actualizar", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Perfil actualizado", user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Usuario no encontrado", err.Error())
		return
	}

	// Filter sensitive fields for public view
	publicUser := models.User{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		Languages: user.Languages,
		CreatedAt: user.CreatedAt,
	}

	utils.SendSuccess(w, http.StatusOK, "Usuario recuperado", publicUser)
}
