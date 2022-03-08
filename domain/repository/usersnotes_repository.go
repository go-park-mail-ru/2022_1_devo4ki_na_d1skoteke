package repository

import "cotion/domain/entity"

type UsersNotesRepository interface {
	AllNotesByUserID(hashedEmail string) ([]entity.Note, error)
	TokensByUserID(hashedEmail string) ([]string, error)
}
