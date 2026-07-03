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

type TourScheduleHandler struct {
	service *service.TourScheduleService
}

func NewTourScheduleHandler(service *service.TourScheduleService) *TourScheduleHandler {
	return &TourScheduleHandler{service: service}
}

func (h *TourScheduleHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tours/{tourId}/schedules", h.ListByTourID)
	mux.HandleFunc("POST /schedules", h.Create)
	mux.HandleFunc("PUT /schedules/{id}", h.Update)
	mux.HandleFunc("DELETE /schedules/{id}", h.Delete)
}

func (h *TourScheduleHandler) ListByTourID(w http.ResponseWriter, r *http.Request) {
	tourIDStr := r.PathValue("tourId")
	tourID, err := strconv.Atoi(tourIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID de tour inválido", err.Error())
		return
	}

	onlyActive := r.URL.Query().Get("active") == "true"

	schedules, err := h.service.ListByTourID(r.Context(), tourID, onlyActive)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al recuperar horarios", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Horarios recuperados", schedules)
}

func (h *TourScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	guideID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	var req models.TourSchedule // Use basic struct to hold raw POST values
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	schedule, err := h.service.Create(r.Context(), guideID, req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al crear horario", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "Horario creado", schedule)
}

func (h *TourScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req models.ScheduleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	schedule, err := h.service.Update(r.Context(), guideID, id, req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al actualizar horario", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Horario actualizado", schedule)
}

func (h *TourScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.Delete(r.Context(), guideID, id)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al eliminar horario", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Horario eliminado", nil)
}
