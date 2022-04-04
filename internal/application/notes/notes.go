package notes

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/infrastructure/storage"
	"cotion/internal/pkg/generator"
	"cotion/internal/pkg/security"
	"errors"
)

type NotesApp struct {
	notesRepository      repository.NotesRepository
	usersNotesRepository repository.UsersNotesRepository
}

func NewNotesApp(notesRepo repository.NotesRepository, usersNotesRepository repository.UsersNotesRepository) *NotesApp {
	return &NotesApp{
		notesRepository:      notesRepo,
		usersNotesRepository: usersNotesRepository,
	}
}

func (n *NotesApp) FindByToken(token string) (entity.Note, error) {
	note, err := n.notesRepository.FindByToken(token)
	if err != nil {
		return entity.Note{}, err
	}
	return note, nil
}

func (n *NotesApp) AllNotesByUserID(hashedEmail string) ([]entity.Note, error) {
	notes, err := n.usersNotesRepository.AllNotesByUserID(hashedEmail)
	if errors.Is(err, storage.CannotFindNotesForUser) {
		return []entity.Note{}, nil
	}
	return notes, err
}

func (n *NotesApp) TokensByUserID(hashedEmail string) ([]string, error) {
	tokens, err := n.usersNotesRepository.TokensByUserID(hashedEmail)
	if errors.Is(err, storage.CannotFindNotesForUser) {
		return []string{}, nil
	}
	return tokens, err
}

func (n *NotesApp) SaveNote(user entity.User, newNote entity.Note) error {
	if err := newNote.Validate(); err != nil {
		return err
	}

	newToken := generator.RandToken()

	if err := n.notesRepository.SaveNote(newToken, newNote); err == nil {
		return err
	}

	if err := n.usersNotesRepository.AddLink(string(security.Hash(user.Email)), newToken); err == nil {
		return err
	}
	return nil
}

func (n *NotesApp) UpdateNote(token string, note entity.Note) error {
	if err := note.Validate(); err != nil {
		return err
	}

	n.notesRepository.UpdateNote(token, note)
	return nil
}

func (n *NotesApp) DeleteNote(userID string, token string) error {
	if err := n.notesRepository.DeleteNote(token); err != nil {
		return err
	}

	if err := n.usersNotesRepository.DeleteLink(userID, token); err != nil {
		return err
	}
	return nil
}
