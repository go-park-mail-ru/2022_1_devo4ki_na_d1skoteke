package repository

import (
	"cotion/internal/api/domain/entity"
)

type UsersNotesRepository interface {
	AddLink(userID string, noteToken string) error
	DeleteLink(userID string, noteToken string) error
	CheckLink(userID string, noteToken string) bool
	AllNotesByUserID(hashedEmail string) ([]entity.Note, error)
}

type NotesRepository interface {
	Save(token string, note entity.Note) error
	Update(token string, note entity.Note) error
	Delete(token string) error
	Find(token string) (entity.Note, error)
}
