package handler

import (
	"cotion/internal/application"
	"cotion/internal/domain/entity"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	userService application.UserAppManager
}

func NewAuthHandler(userService application.UserAppManager) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser entity.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if _, err := h.userService.SaveUser(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
}
