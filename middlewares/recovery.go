package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
)

// RecoveryHandler wraps an HTTP handler and recovers from panics
func RecoveryHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("[PANIC RECOVERED] %s | %s %s | Error: %v\nStack trace:\n%s",
					r.RemoteAddr,
					r.Method,
					r.URL.Path,
					err,
					debug.Stack())

				// Return 500 Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
