// Package calculatorbackend holds small, shared HTTP concerns (currently just
// CORS) that sit above the calculator and handler packages but aren't
// specific to any one binary. cmd/server wires this together with the
// handler package to run the actual service.
package calculatorbackend

import "net/http"

// allowedOrigin is the only origin permitted to call this API: the Vite
// dev server. This is intentionally hardcoded rather than configurable,
// since the project only targets local development at this stage.
const allowedOrigin = "http://localhost:5173"

// WithCORS wraps a handler so that browser requests from allowedOrigin are
// permitted. It also answers CORS preflight (OPTIONS) requests directly,
// since browsers expect those to succeed before sending the real request.
func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
