package application

import (
	"cotion/domain/entity"
	"cotion/domain/repository"
)

type NotesApp struct {
	NotesRepository      repository.NotesRepository
	UsersNotesRepository repository.UsersNotesRepository
}

func NewNotesApp(notesRepo repository.NotesRepository, usersNotesRepository repository.UsersNotesRepository) *NotesApp {
	return &NotesApp{
		NotesRepository:      notesRepo,
		UsersNotesRepository: usersNotesRepository,
	}
}

func (n *NotesApp) GetByToken(token string) (entity.Note, error) {
	note, err := n.NotesRepository.GetByToken(token)
	if err != nil {
		return entity.Note{}, err
	}
	return note, nil
}

func (n *NotesApp) GetAllNotesByUserID(hashedEmail string) ([]entity.Note, error) {
	return n.UsersNotesRepository.GetAllNotesByUserID(hashedEmail)
}
