package handler

import (
	"net/http"
	"strconv"

	"wander/backend/internal/middleware"
	"wander/backend/internal/service"
	"wander/backend/internal/utils"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /payments/bookings/{id}/intent", h.CreateIntent)
	mux.HandleFunc("POST /payments/bookings/{id}/confirm", h.Confirm)
}

func (h *PaymentHandler) CreateIntent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID de reserva inválido", err.Error())
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	resp, err := h.service.CreateIntent(r.Context(), id, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al iniciar el pago", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Intento de pago creado", resp)
}

func (h *PaymentHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID de reserva inválido", err.Error())
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	resp, err := h.service.ConfirmPayment(r.Context(), id, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al confirmar el pago", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Pago confirmado", resp)
}
