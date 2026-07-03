package middleware

import (
	"net/http"
	"strings"

	"wander/backend/internal/utils"
)

// Auth is a simple authentication middleware.
// In production, validate JWT tokens here.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendError(w, http.StatusUnauthorized, "Missing authorization header", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.SendError(w, http.StatusUnauthorized, "Invalid authorization header format", nil)
			return
		}

		token := parts[1]
		_ = token // Validate token in production.

		// TODO: Validate JWT and set user context.

		next.ServeHTTP(w, r)
	})
}
