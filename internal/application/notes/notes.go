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

func (n *NotesApp) SaveNote(user entity.User, newNote entity.Note) error {
	err := newNote.Validate()
	if err != nil {
		return err
	}

	newToken := generator.RandToken()

	err = n.NotesRepository.SaveNote(newToken, newNote)
	if err == nil {
		return err
	}

	err = n.UsersNotesRepository.AddLink(string(security.Hash(user.Email)), newToken)
	if err == nil {
		return err
	}
	return nil
}

func (n *NotesApp) UpdateNote(token string, note entity.Note) error {
	err := note.Validate()
	if err != nil {
		return err
	}

	n.NotesRepository.UpdateNote(token, note)
	return nil
}

func (n *NotesApp) DeleteNote(userID string, token string) error {
	err := n.NotesRepository.DeleteNote(token)
	if err != nil {
		return err
	}
	err = n.UsersNotesRepository.DeleteLink(userID, token)
	if err != nil {
		return err
	}
	return nil
}
