package handler

import (
	"cotion/internal/application"
	"cotion/internal/domain/entity"
	"cotion/internal/pkg/security"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const noteToken = "note-token"

var NoTokenError = errors.New("No token in request.")

type NotesHandler struct {
	notesService  application.NotesAppManager
	authService   application.AuthAppManager
	secureService security.Manager
}

func NewNotesHandler(notesServ application.NotesAppManager, authServ application.AuthAppManager, secureServ security.Manager) *NotesHandler {
	return &NotesHandler{
		notesService:  notesServ,
		authService:   authServ,
		secureService: secureServ,
	}
}

func (h *NotesHandler) ReceiveSingleNote(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "ReceiveSingleNote",
	})

	w.Header().Add("Content-Type", "application/json")
	user := r.Context().Value("user").(entity.User)
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, NoTokenError.Error(), http.StatusBadRequest)
		logger.Warning(NoTokenError)
		return
	}

	userID := h.secureService.Hash(user.Email)
	note, err := h.notesService.GetNote(userID, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(err)
		return
	}
}

func (h *NotesHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "MainPage",
	})

	w.Header().Add("Content-Type", "application/json")

	user := r.Context().Value("user").(entity.User)

	notes, err := h.notesService.AllNotesByUserID(security.Hash(user.Email))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	if err := json.NewEncoder(w).Encode(entity.Notes{Notes: notes}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithFields(log.Fields{
			"notes": entity.Notes{Notes: notes},
		}).Error(err)
		return
	}
}

func (h *NotesHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "CreateNote",
	})

	user := r.Context().Value("user").(entity.User)

	var noteRequest entity.NoteRequest
	if err := noteRequest.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.WithFields(log.Fields{
			"noteRequest": noteRequest,
		}).Warning(err)
		return
	}

	userID := h.secureService.Hash(user.Email)
	if err := h.notesService.SaveNote(userID, noteRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotesHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "UpdateNote",
	})

	user := r.Context().Value("user").(entity.User)
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, NoTokenError.Error(), http.StatusBadRequest)
		logger.Warning(NoTokenError)
		return
	}

	noteRequest := entity.NoteRequest{}
	if err := noteRequest.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Warning(err)
		return
	}

	userID := string(h.secureService.Hash(user.Email))
	if err := h.notesService.UpdateNote(userID, token, noteRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Warning(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotesHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "DeleteNote",
	})

	user := r.Context().Value("user").(entity.User)
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, NoTokenError.Error(), http.StatusBadRequest)
		logger.Warning(NoTokenError)
		return
	}

	userID := h.secureService.Hash(user.Email)
	if err := h.notesService.DeleteNote(userID, token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Warning(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
