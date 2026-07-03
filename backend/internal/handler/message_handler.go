package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"wander/backend/internal/middleware"
	"wander/backend/internal/models"
	"wander/backend/internal/service"
	"wander/backend/internal/utils"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(service *service.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (h *MessageHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /messages/conversations", h.ListConversations)
	mux.HandleFunc("GET /messages/{userId}", h.ListMessages)
	mux.HandleFunc("POST /messages/{userId}", h.Create)
	mux.HandleFunc("GET /messages/stream", h.StreamMessages)
}

func (h *MessageHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	conversations, err := h.service.ListConversations(r.Context(), userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al recuperar conversaciones", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Conversaciones recuperadas", conversations)
}

func (h *MessageHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	otherIDStr := r.PathValue("userId")
	otherID, err := strconv.Atoi(otherIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	messages, err := h.service.ListMessages(r.Context(), userID, otherID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al recuperar mensajes", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Mensajes recuperados", messages)
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	otherIDStr := r.PathValue("userId")
	otherID, err := strconv.Atoi(otherIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	var req models.MessageCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	message, err := h.service.Create(r.Context(), userID, otherID, req.Content, req.BookingID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al enviar mensaje", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "Mensaje enviado", message)
}

func (h *MessageHandler) StreamMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		utils.SendError(w, http.StatusInternalServerError, "Streaming no soportado", nil)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	messageChan := make(chan models.Message, 10)
	h.service.RegisterStream(userID, messageChan)
	defer h.service.UnregisterStream(userID, messageChan)

	slog.Info("SSE Client connected", "user_id", userID)

	// Keep-alive ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-messageChan:
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			_, _ = w.Write([]byte("data: " + string(msgBytes) + "\n\n"))
			flusher.Flush()
		case <-ticker.C:
			// Write keep alive comment
			_, _ = w.Write([]byte(": keepalive\n\n"))
			flusher.Flush()
		case <-r.Context().Done():
			slog.Info("SSE Client disconnected", "user_id", userID)
			return
		}
	}
}
