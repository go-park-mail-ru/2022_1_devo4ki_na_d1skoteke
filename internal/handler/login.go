package handler

import (
	"cotion/internal/application"
	"cotion/internal/domain/entity"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	sessionCookie = "session_id"
)

var ErrDecode = errors.New("problem with decode request")
var ErrNoLoginData = errors.New("no email or password in request")

type LoginHandler struct {
	authService application.AuthAppManager
}

func NewLoginHandler(au application.AuthAppManager) *LoginHandler {
	return &LoginHandler{
		authService: au,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "Login",
	})

	user := entity.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, ErrDecode.Error(), http.StatusBadRequest)
		logger.Warning(err)
		return
	}

	if !user.IsEmail() || !user.IsPassword() {
		http.Error(w, ErrNoLoginData.Error(), http.StatusBadRequest)
		logger.WithFields(log.Fields{
			"userRequest": user,
		}).Warning(ErrNoLoginData)
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
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "Logout",
	})

	sessionCookie, err := r.Cookie(sessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		logger.Error(err)
		return
	}

	newSessionCookie, err := h.authService.Logout(sessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(err)
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
