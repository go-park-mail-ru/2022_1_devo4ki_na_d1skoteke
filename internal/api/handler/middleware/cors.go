package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

func CorsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return enableCORS(next)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://95.163.212.32:3000")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}
