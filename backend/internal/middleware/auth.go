package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"wander/backend/internal/utils"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

// JWTClaims defines the structure of JWT payloads.
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Auth is the JWT authentication middleware.
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			tokenStr := parts[1]
			claims := &JWTClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				utils.SendError(w, http.StatusUnauthorized, "Invalid or expired token", nil)
				return
			}

			// Add properties to request context.
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID helper reads the user ID from the request context.
func GetUserID(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(UserIDKey).(int)
	return val, ok
}

// GetUserRole helper reads the user role from the request context.
func GetUserRole(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(UserRoleKey).(string)
	return val, ok
}

// RequireRole checks if the user has the required role.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := GetUserRole(r.Context())
			if !ok || userRole != role {
				utils.SendError(w, http.StatusForbidden, "Forbidden: insufficient permissions", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
