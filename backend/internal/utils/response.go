package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse writes a JSON response with the given status code.
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Standard API response structure.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SendSuccess sends a success JSON response.
func SendSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	JSONResponse(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends an error JSON response.
func SendError(w http.ResponseWriter, status int, message string, err interface{}) {
	JSONResponse(w, status, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}
