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

type TourHandler struct {
	service *service.TourService
}

func NewTourHandler(service *service.TourService) *TourHandler {
	return &TourHandler{service: service}
}

func (h *TourHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tours", h.List)
	mux.HandleFunc("GET /tours/{id}", h.GetByID)
	mux.HandleFunc("POST /tours", h.Create)
	mux.HandleFunc("PUT /tours/{id}", h.Update)
	mux.HandleFunc("DELETE /tours/{id}", h.Delete)
	mux.HandleFunc("GET /guide/tours", h.ListMyTours)
	mux.HandleFunc("GET /guide/stats", h.GetStats)
}

func (h *TourHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := models.TourFilter{
		Query:        q.Get("q"),
		CategorySlug: q.Get("category"),
		Difficulty:   q.Get("difficulty"),
		Location:     q.Get("location"),
	}

	if catIDStr := q.Get("category_id"); catIDStr != "" {
		if id, err := strconv.Atoi(catIDStr); err == nil {
			filter.CategoryID = id
		}
	}
	if guideIDStr := q.Get("guide_id"); guideIDStr != "" {
		if id, err := strconv.Atoi(guideIDStr); err == nil {
			filter.GuideID = id
		}
	}
	if minPriceStr := q.Get("min_price"); minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &val
		}
	}
	if maxPriceStr := q.Get("max_price"); maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &val
		}
	}
	if limitStr := q.Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = val
		}
	}
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = val
		}
	}

	// Try reading optional user ID for favorites join
	if userID, ok := middleware.GetUserID(r.Context()); ok {
		filter.UserID = userID
	}

	tours, err := h.service.List(r.Context(), filter)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al listar tours", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Tours recuperados", tours)
}

func (h *TourHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	userID := 0
	if uID, ok := middleware.GetUserID(r.Context()); ok {
		userID = uID
	}

	tour, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Tour no encontrado", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Tour recuperado", tour)
}

func (h *TourHandler) Create(w http.ResponseWriter, r *http.Request) {
	guideID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	var req models.TourCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	tour, err := h.service.Create(r.Context(), guideID, req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al crear tour", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "Tour creado exitosamente", tour)
}

func (h *TourHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	guideID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	var req models.TourUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	tour, err := h.service.Update(r.Context(), id, guideID, req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al actualizar tour", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Tour actualizado", tour)
}

func (h *TourHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	guideID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	err = h.service.Delete(r.Context(), id, guideID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al eliminar tour", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Tour eliminado exitosamente", nil)
}

func (h *TourHandler) ListMyTours(w http.ResponseWriter, r *http.Request) {
	guideID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	tours, err := h.service.List(r.Context(), models.TourFilter{GuideID: guideID})
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al recuperar tours", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Tours recuperados", tours)
}

func (h *TourHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	guideID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	stats, err := h.service.GetStats(r.Context(), guideID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al calcular estadísticas", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Estadísticas recuperadas", stats)
}
