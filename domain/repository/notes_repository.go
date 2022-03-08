package repository

import (
	"cotion/domain/entity"
)

type NotesRepository interface {
	FindByToken(token string) (entity.Note, error)
}
