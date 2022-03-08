package application

import (
	"cotion/domain/entity"
	"net/http"
)

type NotesAppManager interface {
	GetByToken(token string) (entity.Note, error)
	GetAllNotesByUserID(hashedEmail string) ([]entity.Note, error)
}

type UserAppManager interface {
	SaveUser(registerUser entity.RegisterUserRequest) (entity.User, error)
	GetUser(email string) (entity.User, error)
}

type AuthAppManager interface {
	Login(login string, password string) (*http.Cookie, error)
	Logout(sessionCookie *http.Cookie) (*http.Cookie, error)
	Auth(SID string) bool
}
