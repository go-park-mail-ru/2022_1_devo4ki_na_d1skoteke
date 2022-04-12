package handler

import (
	"cotion/internal/application"
	"cotion/internal/domain/entity"
	"cotion/internal/pkg/security"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const packageName = "handler"

type UserHandler struct {
	userService application.UserAppManager
}

func NewUserHandler(userService application.UserAppManager) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "SignUp",
	})
	var newUser entity.UserRequest
	if err := newUser.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Warning(err)
		return
	}

	if err := h.userService.Save(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Warning(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entity.User)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "GetUser",
			"user":     user,
		}).Error(err)
		return
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "UpdateUser",
	})

	curUser := r.Context().Value("user").(entity.User)

	var updateUser entity.UserRequest
	if err := updateUser.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Warning(err)
		return
	}

	if err := h.userService.Update(curUser, updateUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entity.User)

	if err := h.userService.Delete(security.Hash(user.Email)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Delete",
		}).Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
