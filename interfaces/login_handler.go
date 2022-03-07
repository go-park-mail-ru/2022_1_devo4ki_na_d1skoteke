package interfaces

import (
	"cotion/application"
	"cotion/domain/entity"
	"encoding/json"
	"net/http"
)

const (
	emailJson = "email"
	passJson  = "password"
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
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "problem with decode request", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
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
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	newSessionCookie, err := h.authService.Logout(sessionCookie)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, newSessionCookie)
	w.WriteHeader(http.StatusOK)
}

func (h *LoginHandler) Auth(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	isAuth := h.authService.Auth(sessionCookie.Value)

	if !isAuth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
