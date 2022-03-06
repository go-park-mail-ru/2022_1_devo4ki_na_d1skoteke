package application

import (
	"cotion/domain/entity"
	"cotion/domain/repository"
)

type NotesApp struct {
	NotesRepository repository.NotesRepository
}

func NewNotesApp(notesRepo repository.NotesRepository) *NotesApp {
	return &NotesApp{
		NotesRepository: notesRepo,
	}
}

func (n *NotesApp) GetByToken(token string) (entity.Note, error) {
	note, err := n.NotesRepository.GetByToken(token)
	if err != nil {
		return entity.Note{}, err
	}
	return note, nil
}
