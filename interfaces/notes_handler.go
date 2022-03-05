package interfaces

import (
	"cotion/application"
	"cotion/infrastructure/security"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const noteToken = "note-token"

type NotesHandler struct {
	notesService application.NotesAppManager
}

func NewNotesHandler(service application.NotesAppManager) *NotesHandler {
	return &NotesHandler{
		notesService: service,
	}
}

func (h *NotesHandler) ReceiveSingleNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, ok := vars[noteToken]
	if !ok {
		http.Error(w, "no token in request", http.StatusBadRequest)
		return
	}

	note, err := h.notesService.GetByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *NotesHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	notes, err := h.notesService.GetAllNotesByUserID(string(security.Hash("email@vk.team")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
