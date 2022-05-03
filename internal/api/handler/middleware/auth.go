package middleware

import (
	"context"
	"cotion/internal/api/application"
	"errors"
	"net/http"
)

const (
	sessionCookie = "session_id"
)

var ErrUnauthorized = errors.New("user is not authorized")
var ErrAuthorized = errors.New("user is already authorized")

type AuthMiddleware struct {
	authService application.AuthAppManager
}

func NewAuthMiddleware(authServ application.AuthAppManager) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authServ,
	}
}

func (amw *AuthMiddleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sCookie, err := r.Cookie(sessionCookie)
		if err != nil {
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		user, ok := amw.authService.Auth(sCookie)
		if !ok {
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
	}
}

func (amw *AuthMiddleware) NotAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sCookie, err := r.Cookie(sessionCookie)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		_, ok := amw.authService.Auth(sCookie)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, ErrAuthorized.Error(), http.StatusBadRequest)
	}
}
