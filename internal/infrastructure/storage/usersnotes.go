package storage

import (
	"cotion/internal/domain/entity"
	"cotion/internal/pkg/security"
	"errors"
	"sync"
)

type UsersNotesStorage struct {
	data  sync.Map
	notes *NotesStorage
}

var CannotFindNotesForUser = errors.New("cannot find notes")
var CannotFindNoteByToken = errors.New("cannot find note by token")

func NewUsersNotesStorage(notesStorage *NotesStorage) *UsersNotesStorage {
	storage := &UsersNotesStorage{
		data:  sync.Map{},
		notes: notesStorage,
	}
	storage.data.Store(string(security.Hash("test@mail.ru")), []string{"1", "3"})
	storage.data.Store(string(security.Hash("test3@mail.ru")), []string{"0"})
	storage.data.Store(string(security.Hash("nikita@mail.ru")), []string{"2"})
	return storage
}

func (storage *UsersNotesStorage) AllNotesByUserID(hashedEmail string) ([]entity.Note, error) {
	rawNotesIDs, ok := storage.data.Load(hashedEmail)
	if !ok {
		return []entity.Note{}, CannotFindNotesForUser
	}

	notesIDs := rawNotesIDs.([]string)
	notes := make([]entity.Note, 0)

	for _, id := range notesIDs {
		note, err := storage.notes.FindByToken(id)
		if err != nil {
			return []entity.Note{}, CannotFindNoteByToken
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (storage *UsersNotesStorage) TokensByUserID(hashedEmail string) ([]string, error) {
	rawNotesIDs, ok := storage.data.Load(hashedEmail)
	if !ok {
		return []string{}, CannotFindNotesForUser
	}
	return rawNotesIDs.([]string), nil
}

func (storage *UsersNotesStorage) AddLink(userID string, noteToken string) error {
	rawNotesIDs, ok := storage.data.Load(userID)
	if !ok {
		storage.data.Store(userID, []string{noteToken})
	}
	NotesIDs := rawNotesIDs.([]string)
	NotesIDs = append(NotesIDs, noteToken)
	storage.data.Store(userID, NotesIDs)
	return nil
}

func findNote(NotesIDs []string, token string) (int, bool) {
	for index, curToken := range NotesIDs {
		if curToken == token {
			return index, true
		}
	}
	return -1, false
}

func (storage *UsersNotesStorage) DeleteLink(userID string, noteToken string) error {
	rawNotesIDs, ok := storage.data.Load(userID)
	if !ok {
		return errors.New("can't find user")
	}
	NotesIDs := rawNotesIDs.([]string)

	noteIndex, ok := findNote(NotesIDs, noteToken)
	if !ok {
		return errors.New("can't find token in user's notes")
	}

	NotesIDs[noteIndex] = NotesIDs[len(NotesIDs)-1]
	NotesIDs = NotesIDs[:len(NotesIDs)-1]
	storage.data.Store(userID, NotesIDs)
	return nil
}

func (storage *UsersNotesStorage) CheckLink(userID string, noteToken string) bool {
	rawNotesIDs, ok := storage.data.Load(userID)
	if !ok {
		return false
	}
	NotesIDs := rawNotesIDs.([]string)
	_, ok = findNote(NotesIDs, noteToken)
	return ok
}
