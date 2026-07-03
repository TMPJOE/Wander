package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"wander/backend/internal/utils"
)

// Recovery recovers from panics and returns a 500 response, logging details using slog.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := string(debug.Stack())
				slog.Error("panic recovered",
					slog.Any("recover", rec),
					slog.String("stack", stack),
					slog.String("path", r.URL.Path),
				)
				utils.SendError(w, http.StatusInternalServerError, "Internal Server Error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
