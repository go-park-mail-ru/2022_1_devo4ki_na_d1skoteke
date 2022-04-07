package handler

import (
	"cotion/internal/application"
	"cotion/internal/domain/entity"
	"encoding/json"
	"net/http"
)

const (
	sessionCookie = "session_id"
)

type LoginHandler struct {
	authService application.AuthAppManager
}

func NewLoginHandler(au application.AuthAppManager) *LoginHandler {
	return &LoginHandler{
		authService: au,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	user := entity.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "problem with decode request", http.StatusBadRequest)
		return
	}

	if !user.IsEmail() || !user.IsPassword() {
		http.Error(w, "no email or password in request", http.StatusBadRequest)
		return
	}

	cookie, err := h.authService.Login(user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(sessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	newSessionCookie, err := h.authService.Logout(sessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, newSessionCookie)
	w.WriteHeader(http.StatusOK)
}

func (h *LoginHandler) Auth(w http.ResponseWriter, r *http.Request) {
	sCookie, err := r.Cookie(sessionCookie)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, auth := h.authService.Auth(sCookie); !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
