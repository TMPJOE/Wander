package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"wander/backend/internal/utils"
)

// Recovery recovers from panics and returns a 500 response.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := debug.Stack()
				fmt.Printf("[PANIC] %v\n%s\n", rec, stack)
				utils.SendError(w, http.StatusInternalServerError, "Internal Server Error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
