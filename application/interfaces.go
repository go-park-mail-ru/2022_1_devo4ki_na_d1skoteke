package application

import (
	"cotion/domain/entity"
)

type NotesAppManager interface {
	GetByToken(token string) (entity.Note, error)
	GetAllNotesByUserID(hashedEmail string) ([]entity.Note, error)
}

type UserAppManager interface {
	SaveUser(registerUser entity.RegisterUserRequest) (entity.User, error)
	GetUser(email string) (entity.User, error)
}
