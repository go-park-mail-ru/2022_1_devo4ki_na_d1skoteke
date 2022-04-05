package handler

import (
	"cotion/internal/application"
	"cotion/internal/domain/entity"
	"cotion/internal/pkg/security"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
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
	w.Header().Add("Content-Type", "application/json")
	user := r.Context().Value("user").(entity.User)
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, "no token in request", http.StatusBadRequest)
		return
	}

	note, err := h.notesService.GetNote(user, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NotesHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	user := r.Context().Value("user").(entity.User)

	notes, err := h.notesService.AllNotesByUserID(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(entity.Notes{Notes: notes}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NotesHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entity.User)

	var noteRequest entity.NoteRequest
	if err := noteRequest.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.notesService.SaveNote(user, noteRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotesHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entity.User)
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, "no token in request", http.StatusBadRequest)
		return
	}

	noteRequest := entity.NoteRequest{}
	if err := noteRequest.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.notesService.UpdateNote(user, token, noteRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotesHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entity.User)
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, NoTokenError.Error(), http.StatusBadRequest)
		return
	}

	if err := h.notesService.DeleteNote(user, token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
