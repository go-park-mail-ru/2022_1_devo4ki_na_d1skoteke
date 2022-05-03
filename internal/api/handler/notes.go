package handler

import (
	"context"
	"cotion/internal/api/application"
	"cotion/internal/api/domain/entity"
	grpcNoteServ "cotion/internal/note/infra/grpc"
	"cotion/internal/pkg/security"
	"cotion/internal/pkg/xss"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const noteToken = "note-token"

var NoTokenError = errors.New("No token in request.")

type NotesHandler struct {
	notesService  grpcNoteServ.NoteServiceClient
	authService   application.AuthAppManager
	secureService security.Manager
}

func NewNotesHandler(notesServ grpcNoteServ.NoteServiceClient, authServ application.AuthAppManager, secureServ security.Manager) *NotesHandler {
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

	ctx := context.Background()
	grpcNoteReq := &grpcNoteServ.NoteReq{
		UserID: user.UserID,
		NoteID: token,
	}
	grpcNote, err := h.notesService.Get(ctx, grpcNoteReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	note := entity.Note{
		Name: grpcNote.Name,
		Body: grpcNote.Body,
	}

	xss.SanitizeNote(&note)

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
	ctx := context.Background()
	grpcNoteReq := &grpcNoteServ.NoteReq{
		UserID: user.UserID,
	}

	notes := make([]entity.Note, 0)
	grpcStream, err := h.notesService.NotesList(ctx, grpcNoteReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	for {
		grpcNote, err := grpcStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error(err)
			return
		}
		note := entity.Note{
			Name: grpcNote.Name,
			Body: grpcNote.Body,
		}
		notes = append(notes, note)
	}

	xss.SanitizeNotes(&notes)

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

	ctx := context.Background()
	grpcNoteReq := &grpcNoteServ.NoteReq{
		UserID: user.UserID,
		Note: &grpcNoteServ.Note{
			Name: noteRequest.Name,
			Body: noteRequest.Body,
		},
	}
	if _, err := h.notesService.Save(ctx, grpcNoteReq); err != nil {
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

	ctx := context.Background()
	grpcNoteReq := &grpcNoteServ.NoteReq{
		UserID: user.UserID,
		NoteID: token,
		Note: &grpcNoteServ.Note{
			Name: noteRequest.Name,
			Body: noteRequest.Body,
		},
	}
	if _, err := h.notesService.Update(ctx, grpcNoteReq); err != nil {
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

	ctx := context.Background()
	grpcNoteReq := &grpcNoteServ.NoteReq{
		UserID: user.UserID,
		NoteID: token,
	}
	if _, err := h.notesService.Delete(ctx, grpcNoteReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Warning(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
