package middleware

import (
	"net/http"
)

// Json middleware sets the response's content type to application/json
func Json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
