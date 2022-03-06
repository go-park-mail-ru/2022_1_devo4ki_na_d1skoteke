package interfaces

import (
	"cotion/application"
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

	json.NewEncoder(w).Encode(note)
}
