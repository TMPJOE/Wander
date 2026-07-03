package handler

import (
	"encoding/json"
	"net/http"

	"wander/backend/internal/models"
	"wander/backend/internal/service"
	"wander/backend/internal/utils"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	res, err := h.service.Register(r.Context(), req)
	if err != nil {
		utils.SendError(w, http.StatusConflict, "Error al registrar usuario", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "Usuario registrado exitosamente", res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Request incorrecto", err.Error())
		return
	}

	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Credenciales incorrectas", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Login exitoso", res)
}
