package application

import (
	"cotion/internal/domain/entity"
	"net/http"
)

type NotesAppManager interface {
	FindByToken(token string) (entity.Note, error)
	AllNotesByUserID(user entity.User) ([]entity.Note, error)
	TokensByUserID(hashedEmail string) ([]string, error)
	SaveNote(user entity.User, noteRequest entity.NoteRequest) error
	GetNote(user entity.User, token string) (entity.Note, error)
	UpdateNote(user entity.User, token string, noteRequest entity.NoteRequest) error
	DeleteNote(user entity.User, token string) error
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
