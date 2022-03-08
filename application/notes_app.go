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

func (n *NotesApp) FindByToken(token string) (entity.Note, error) {
	note, err := n.NotesRepository.FindByToken(token)
	if err != nil {
		return entity.Note{}, err
	}
	return note, nil
}

func (n *NotesApp) AllNotesByUserID(hashedEmail string) ([]entity.Note, error) {
	return n.UsersNotesRepository.AllNotesByUserID(hashedEmail)
}

func (n *NotesApp) TokensByUserID(hashedEmail string) ([]string, error) {
	return n.UsersNotesRepository.TokensByUserID(hashedEmail)
}
