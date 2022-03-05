package repository

import "cotion/domain/entity"

type UsersNotesRepository interface {
	GetAllNotesByUserID(hashedEmail string) ([]entity.Note, error)
}
