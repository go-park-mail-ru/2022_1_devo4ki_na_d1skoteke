package middleware

import (
	"context"
	"cotion/internal/application"
	"errors"
	"net/http"
)

const (
	sessionCookie = "session_id"
)

var UnauthorizedError = errors.New("user is not authorized")
var AuthorizedError = errors.New("user is already authorized")

type AuthMiddleware struct {
	authService application.AuthAppManager
}

func NewAuthMiddleware(authServ application.AuthAppManager) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authServ,
	}
}

func (amw *AuthMiddleware) AuthMiddleware(authType bool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch authType {
		case true:
			sCookie, err := r.Cookie(sessionCookie)
			if err != nil {
				http.Error(w, UnauthorizedError.Error(), http.StatusUnauthorized)
				return
			}

			user, ok := amw.authService.Auth(sCookie)
			if !ok {
				http.Error(w, UnauthorizedError.Error(), http.StatusUnauthorized)
				return
			}
			
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))

		case false:
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

			http.Error(w, AuthorizedError.Error(), http.StatusBadRequest)
		}
	}
}
