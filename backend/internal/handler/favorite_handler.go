package handler

import (
	"net/http"
	"strconv"

	"wander/backend/internal/middleware"
	"wander/backend/internal/service"
	"wander/backend/internal/utils"
)

type FavoriteHandler struct {
	service *service.FavoriteService
}

func NewFavoriteHandler(service *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

func (h *FavoriteHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /favorites", h.List)
	mux.HandleFunc("POST /favorites/{tourId}", h.Add)
	mux.HandleFunc("DELETE /favorites/{tourId}", h.Remove)
}

func (h *FavoriteHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	tours, err := h.service.List(r.Context(), userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al listar favoritos", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Favoritos recuperados", tours)
}

func (h *FavoriteHandler) Add(w http.ResponseWriter, r *http.Request) {
	tourIDStr := r.PathValue("tourId")
	tourID, err := strconv.Atoi(tourIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID de tour inválido", err.Error())
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	err = h.service.Add(r.Context(), userID, tourID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al añadir favorito", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Añadido a favoritos", nil)
}

func (h *FavoriteHandler) Remove(w http.ResponseWriter, r *http.Request) {
	tourIDStr := r.PathValue("tourId")
	tourID, err := strconv.Atoi(tourIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID de tour inválido", err.Error())
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	err = h.service.Remove(r.Context(), userID, tourID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al quitar favorito", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Quitado de favoritos", nil)
}
