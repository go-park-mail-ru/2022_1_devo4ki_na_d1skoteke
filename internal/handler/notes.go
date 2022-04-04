package handler

import (
	"cotion/internal/application"
	"cotion/internal/domain/entity"
	"cotion/internal/pkg/contains"
	"cotion/internal/pkg/security"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const noteToken = "note-token"

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
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, "no token in request", http.StatusBadRequest)
		return
	}

	user, auth := isAuth(h.authService, r)
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	userTokens, err := h.notesService.TokensByUserID(string(h.secureService.Hash(user.Email)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !contains.Contains(userTokens, token) {
		http.Error(w, "permission deny", http.StatusMethodNotAllowed)
		return
	}

	note, err := h.notesService.FindByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *NotesHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	user, auth := isAuth(h.authService, r)
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	notes, err := h.notesService.AllNotesByUserID(string(security.Hash(user.Email)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(entity.Notes{Notes: notes}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *NotesHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	user, auth := isAuth(h.authService, r)
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var newNote entity.Note
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := h.notesService.SaveNote(user, newNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotesHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, "no token in request", http.StatusBadRequest)
		return
	}

	user, auth := isAuth(h.authService, r)
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userTokens, err := h.notesService.TokensByUserID(string(h.secureService.Hash(user.Email)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !contains.Contains(userTokens, token) {
		http.Error(w, "permission deny", http.StatusMethodNotAllowed)
		return
	}

	var newNote entity.Note
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := h.notesService.UpdateNote(token, newNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotesHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, "no token in request", http.StatusBadRequest)
		return
	}

	user, auth := isAuth(h.authService, r)
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userTokens, err := h.notesService.TokensByUserID(string(h.secureService.Hash(user.Email)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !contains.Contains(userTokens, token) {
		http.Error(w, "permission deny", http.StatusMethodNotAllowed)
		return
	}

	if err := h.notesService.DeleteNote(string(security.Hash(user.Email)), token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
