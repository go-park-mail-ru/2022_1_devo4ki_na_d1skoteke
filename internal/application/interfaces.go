package application

import (
	"cotion/internal/domain/entity"
	"net/http"
)

type NotesAppManager interface {
	FindByToken(token string) (entity.Note, error)
	AllNotesByUserID(userID string) ([]entity.Note, error)
	TokensByUserID(hashedEmail string) ([]string, error)
	SaveNote(userID string, noteRequest entity.NoteRequest) error
	GetNote(userID string, noteToken string) (entity.Note, error)
	UpdateNote(userID string, noteToken string, noteRequest entity.NoteRequest) error
	DeleteNote(userID string, noteToken string) error
}

type UserAppManager interface {
	SaveUser(registerUser entity.RegisterUserRequest) (entity.User, error)
	GetUser(email string) (entity.User, error)
}

type AuthAppManager interface {
	Login(login string, password string) (*http.Cookie, error)
	Logout(sessionCookie *http.Cookie) (*http.Cookie, error)
	Auth(sessionCookie *http.Cookie) (entity.User, bool)
}
