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

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /bookings", h.Create)
	mux.HandleFunc("GET /bookings", h.List)
	mux.HandleFunc("GET /bookings/{id}", h.GetByID)
	mux.HandleFunc("PATCH /bookings/{id}/cancel", h.Cancel)
	mux.HandleFunc("PATCH /bookings/{id}/confirm", h.Confirm)
	mux.HandleFunc("PATCH /bookings/{id}/complete", h.Complete)
	mux.HandleFunc("PATCH /bookings/{id}/reject", h.Reject)
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	var req models.BookingCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	booking, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al crear reserva", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "Reserva creada", booking)
}

func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	role, ok2 := middleware.GetUserRole(r.Context())
	if !ok || !ok2 {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	booking, err := h.service.GetByID(r.Context(), id, userID, role)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al recuperar reserva", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Reserva recuperada", booking)
}

func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	role, ok2 := middleware.GetUserRole(r.Context())
	if !ok || !ok2 {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	bookings, err := h.service.ListByUser(r.Context(), userID, role)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al listar reservas", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Reservas recuperadas", bookings)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	err = h.service.Cancel(r.Context(), id, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al cancelar", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Reserva cancelada exitosamente", nil)
}

func (h *BookingHandler) Confirm(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.Confirm(r.Context(), id, guideID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al confirmar", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Reserva confirmada exitosamente", nil)
}

func (h *BookingHandler) Complete(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.Complete(r.Context(), id, guideID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al marcar completada", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Reserva marcada como completada", nil)
}

func (h *BookingHandler) Reject(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.Reject(r.Context(), id, guideID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al rechazar", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Reserva rechazada exitosamente", nil)
}
