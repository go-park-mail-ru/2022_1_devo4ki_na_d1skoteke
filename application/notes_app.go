package application

import (
	"cotion/domain/entity"
	"cotion/domain/repository"
	"cotion/infrastructure/storage"
	"errors"
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
	notes, err := n.UsersNotesRepository.AllNotesByUserID(hashedEmail)
	if errors.Is(err, storage.CannotFindNotesForUser) {
		return []entity.Note{}, nil
	}
	return notes, err
}

func (n *NotesApp) TokensByUserID(hashedEmail string) ([]string, error) {
	tokens, err := n.UsersNotesRepository.TokensByUserID(hashedEmail)
	if errors.Is(err, storage.CannotFindNotesForUser) {
		return []string{}, nil
	}
	return tokens, err
}
