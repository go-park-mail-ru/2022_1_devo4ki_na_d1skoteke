package notes

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/pkg/generator"
	"errors"
	log "github.com/sirupsen/logrus"
)

const packageName = "app notes"

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

func (n *NotesApp) AllNotesByUserID(userID string) ([]entity.Note, error) {
	return n.usersNotesRepository.AllNotesByUserID(userID)
}

func (n *NotesApp) SaveNote(userID string, noteRequest entity.NoteRequest) error {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "SaveNote",
	})

	newToken := generator.RandToken()
	newNote := entity.Note{
		Name: noteRequest.Name,
		Body: noteRequest.Body,
	}

	if err := n.notesRepository.Save(newToken, newNote); err != nil {
		logger.Error(err)
		return err
	}

	if err := n.usersNotesRepository.AddLink(userID, newToken); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (n *NotesApp) GetNote(userID string, noteToken string) (entity.Note, error) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "GetNote",
	})

	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		logger.Warning(ErrNoteAccess)
		return entity.Note{}, ErrNoteAccess
	}

	note, err := n.notesRepository.Find(noteToken)
	if err != nil {
		logger.Error(err)
		return entity.Note{}, err
	}

	return note, nil
}

func (n *NotesApp) UpdateNote(userID string, noteToken string, noteRequest entity.NoteRequest) error {
	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "UpdateNote",
		}).Warning(ErrNoteAccess)
		return ErrNoteAccess
	}

	updateNote := entity.Note{
		Name: noteRequest.Name,
		Body: noteRequest.Body,
	}

	return n.notesRepository.Update(noteToken, updateNote)
}

func (n *NotesApp) DeleteNote(userID string, noteToken string) error {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "DeleteNote",
	})

	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		logger.Warning(ErrNoteAccess)
		return ErrNoteAccess
	}

	if err := n.notesRepository.Delete(noteToken); err != nil {
		logger.Error(err)
		return err
	}

	return n.usersNotesRepository.DeleteLink(userID, noteToken)
}
