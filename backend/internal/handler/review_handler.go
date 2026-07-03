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

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(service *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tours/{tourId}/reviews", h.ListByTour)
	mux.HandleFunc("POST /tours/{tourId}/reviews", h.Create)
}

func (h *ReviewHandler) ListByTour(w http.ResponseWriter, r *http.Request) {
	tourIDStr := r.PathValue("tourId")
	tourID, err := strconv.Atoi(tourIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID de tour inválido", err.Error())
		return
	}

	reviews, err := h.service.ListByTourID(r.Context(), tourID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al recuperar reseñas", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Reseñas recuperadas", reviews)
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req models.ReviewCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	review, err := h.service.Create(r.Context(), userID, tourID, req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al publicar reseña", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "Reseña publicada", review)
}
