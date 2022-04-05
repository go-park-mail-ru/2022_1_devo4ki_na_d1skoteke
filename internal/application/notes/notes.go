package notes

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/infrastructure/storage"
	"cotion/internal/pkg/generator"
	"cotion/internal/pkg/security"
	"errors"
)

var NoteAccessError = errors.New("The user does not have access to this note. Or the note does not exist.")

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

func (n *NotesApp) AllNotesByUserID(user entity.User) ([]entity.Note, error) {
	notes, err := n.usersNotesRepository.AllNotesByUserID(string(security.Hash(user.Email)))
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

func (n *NotesApp) SaveNote(user entity.User, noteRequest entity.NoteRequest) error {
	newToken := generator.RandToken()

	newNote := entity.Note{
		Name: noteRequest.Name,
		Body: noteRequest.Body,
	}

	if err := n.notesRepository.SaveNote(newToken, newNote); err == nil {
		return err
	}

	if err := n.usersNotesRepository.AddLink(string(security.Hash(user.Email)), newToken); err == nil {
		return err
	}
	return nil
}

func (n *NotesApp) GetNote(user entity.User, noteToken string) (entity.Note, error) {
	if !n.usersNotesRepository.CheckLink(string(security.Hash(user.Email)), noteToken) {
		return entity.Note{}, NoteAccessError
	}

	note, err := n.notesRepository.FindByToken(noteToken)
	if err != nil {
		return entity.Note{}, err
	}

	return note, nil
}

func (n *NotesApp) UpdateNote(user entity.User, noteToken string, noteRequest entity.NoteRequest) error {
	if !n.usersNotesRepository.CheckLink(string(security.Hash(user.Email)), noteToken) {
		return NoteAccessError
	}

	updateNote := entity.Note{
		Name: noteRequest.Name,
		Body: noteRequest.Body,
	}

	n.notesRepository.UpdateNote(noteToken, updateNote)
	return nil
}

func (n *NotesApp) DeleteNote(user entity.User, noteToken string) error {
	userID := string(security.Hash(user.Email))
	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		return NoteAccessError
	}

	if err := n.notesRepository.DeleteNote(noteToken); err != nil {
		return err
	}

	if err := n.usersNotesRepository.DeleteLink(userID, noteToken); err != nil {
		return err
	}
	return nil
}
