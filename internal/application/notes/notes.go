package notes

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/infrastructure/storage"
	"cotion/internal/pkg/generator"
	"errors"
)

var ErrNoteAccess = errors.New("The user does not have access to this note. Or the note does not exist.")

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
	note, err := n.notesRepository.Find(token)
	if err != nil {
		return entity.Note{}, err
	}
	return note, nil
}

func (n *NotesApp) AllNotesByUserID(userID string) ([]entity.Note, error) {
	notes, err := n.usersNotesRepository.AllNotesByUserID(userID)
	if errors.Is(err, storage.ErrFindNotesForUser) {
		return []entity.Note{}, nil
	}
	return notes, err
}

func (n *NotesApp) SaveNote(userID string, noteRequest entity.NoteRequest) error {
	newToken := generator.RandToken()

	newNote := entity.Note{
		Name: noteRequest.Name,
		Body: noteRequest.Body,
	}

	if err := n.notesRepository.Save(newToken, newNote); err != nil {
		return err
	}

	if err := n.usersNotesRepository.AddLink(userID, newToken); err != nil {
		return err
	}
	return nil
}

func (n *NotesApp) GetNote(userID string, noteToken string) (entity.Note, error) {
	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		return entity.Note{}, ErrNoteAccess
	}

	note, err := n.notesRepository.Find(noteToken)
	if err != nil {
		return entity.Note{}, err
	}

	return note, nil
}

func (n *NotesApp) UpdateNote(userID string, noteToken string, noteRequest entity.NoteRequest) error {
	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		return ErrNoteAccess
	}

	updateNote := entity.Note{
		Name: noteRequest.Name,
		Body: noteRequest.Body,
	}

	n.notesRepository.Update(noteToken, updateNote)
	return nil
}

func (n *NotesApp) DeleteNote(userID string, noteToken string) error {
	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		return ErrNoteAccess
	}

	if err := n.notesRepository.Delete(noteToken); err != nil {
		return err
	}

	if err := n.usersNotesRepository.DeleteLink(userID, noteToken); err != nil {
		return err
	}
	return nil
}
