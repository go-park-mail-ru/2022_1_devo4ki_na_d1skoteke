package application

import (
	"cotion/domain/entity"
	"net/http"
)

type NotesAppManager interface {
	FindByToken(token string) (entity.Note, error)
	AllNotesByUserID(hashedEmail string) ([]entity.Note, error)
	TokensByUserID(hashedEmail string) ([]string, error)
}

type UserAppManager interface {
	SaveUser(registerUser entity.RegisterUserRequest) (entity.User, error)
	GetUser(email string) (entity.User, error)
}

type AuthAppManager interface {
	Login(login string, password string) (*http.Cookie, error)
	Logout(sessionCookie http.Cookie) (*http.Cookie, error)
	Auth(r *http.Request) (entity.User, bool)
}
