package repository

import "cotion/internal/domain/entity"

type SessionRepository interface {
	HasSession(SID string) (entity.Session, bool)
	NewSession(SID string, user entity.User) (entity.Session, error)
	DeleteSession(SID string)
}

type UserRepository interface {
	SaveUser(user entity.User) (entity.User, error)
	GetUser(email string) (entity.User, error)
}

type UsersNotesRepository interface {
	AllNotesByUserID(hashedEmail string) ([]entity.Note, error)
	TokensByUserID(hashedEmail string) ([]string, error)
}

type NotesRepository interface {
	FindByToken(token string) (entity.Note, error)
}