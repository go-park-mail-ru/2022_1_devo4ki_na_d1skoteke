package interfaces

import (
	"cotion/application"
	"cotion/domain/entity"
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
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	_, err = h.userService.SaveUser(newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
}
