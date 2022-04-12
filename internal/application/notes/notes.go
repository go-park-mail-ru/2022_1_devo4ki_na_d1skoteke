package notes

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/pkg/generator"
	"errors"
	log "github.com/sirupsen/logrus"
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

func (n *NotesApp) AllNotesByUserID(userID string) ([]entity.Note, error) {
	return n.usersNotesRepository.AllNotesByUserID(userID)
}

func (n *NotesApp) SaveNote(userID string, noteRequest entity.NoteRequest) error {
	newToken := generator.RandToken()

	newNote := entity.Note{
		Name: noteRequest.Name,
		Body: noteRequest.Body,
	}

	if err := n.notesRepository.Save(newToken, newNote); err != nil {
		log.WithFields(log.Fields{
			"package":  "application note",
			"function": "SaveNote",
		}).Error(err)
		return err
	}

	if err := n.usersNotesRepository.AddLink(userID, newToken); err != nil {
		log.WithFields(log.Fields{
			"package":  "application note",
			"function": "SaveNote",
		}).Error(err)
		return err
	}
	return nil
}

func (n *NotesApp) GetNote(userID string, noteToken string) (entity.Note, error) {
	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		log.WithFields(log.Fields{
			"package":  "application note",
			"function": "GetNote",
		}).Warning(ErrNoteAccess)
		return entity.Note{}, ErrNoteAccess
	}

	note, err := n.notesRepository.Find(noteToken)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "application note",
			"function": "GetNote",
		}).Error(err)
		return entity.Note{}, err
	}

	return note, nil
}

func (n *NotesApp) UpdateNote(userID string, noteToken string, noteRequest entity.NoteRequest) error {
	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		log.WithFields(log.Fields{
			"package":  "application note",
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
	if !n.usersNotesRepository.CheckLink(userID, noteToken) {
		log.WithFields(log.Fields{
			"package":  "application note",
			"function": "DeleteNote",
		}).Warning(ErrNoteAccess)
		return ErrNoteAccess
	}

	if err := n.notesRepository.Delete(noteToken); err != nil {
		log.WithFields(log.Fields{
			"package":  "application note",
			"function": "DeleteNote",
		}).Error(err)
		return err
	}

	return n.usersNotesRepository.DeleteLink(userID, noteToken)
}
