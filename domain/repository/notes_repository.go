package repository

import (
	"cotion/domain/entity"
)

type NotesRepository interface {
	GetByToken(token string) (entity.Note, error)
}
