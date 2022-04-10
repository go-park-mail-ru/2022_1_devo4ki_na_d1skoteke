package handler

import (
    "cotion/internal/application"
    "cotion/internal/domain/entity"
    "encoding/json"
    "net/http"
)

type UserHandler struct {
    userService application.UserAppManager
}

func NewUserHandler(userService application.UserAppManager) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
    var newUser entity.RegisterUserRequest
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    if _, err := h.userService.Save(newUser); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    var user entity.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    if _, err := h.userService.Update(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    var user entity.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    if err := h.userService.Delete(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusOK)
}
