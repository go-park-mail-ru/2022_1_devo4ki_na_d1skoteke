package storage

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/pkg/security"
	"errors"
	"sync"
)

type UsersNotesStorage struct {
	data  sync.Map
	notes repository.NotesRepository
}

var ErrFindNotesForUser = errors.New("cannot find notes")
var ErrFindNoteByToken = errors.New("cannot find note by token")
var ErrFindUser = errors.New("can't find user")
var ErrFindTokenInUsersNotes = errors.New("can't find token in user's notes")

func NewUsersNotesStorage(notesStorage repository.NotesRepository) *UsersNotesStorage {
	storage := &UsersNotesStorage{
		data:  sync.Map{},
		notes: notesStorage,
	}
	storage.data.Store(security.Hash("test@mail.ru"), []string{"1", "3"})
	storage.data.Store(security.Hash("test3@mail.ru"), []string{"0"})
	storage.data.Store(security.Hash("nikita@mail.ru"), []string{"2"})
	return storage
}

func (storage *UsersNotesStorage) AllNotesByUserID(hashedEmail string) ([]entity.Note, error) {
	rawNotesIDs, ok := storage.data.Load(hashedEmail)
	if !ok {
		return []entity.Note{}, ErrFindNotesForUser
	}

	notesIDs := rawNotesIDs.([]string)
	notes := make([]entity.Note, 0)

	for _, id := range notesIDs {
		note, err := storage.notes.Find(id)
		if err != nil {
			return []entity.Note{}, ErrFindNoteByToken
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (storage *UsersNotesStorage) TokensByUserID(hashedEmail string) ([]string, error) {
	rawNotesIDs, ok := storage.data.Load(hashedEmail)
	if !ok {
		return []string{}, ErrFindNotesForUser
	}
	return rawNotesIDs.([]string), nil
}

func (storage *UsersNotesStorage) AddLink(userID string, noteToken string) error {
	rawNotesIDs, ok := storage.data.Load(userID)
	if !ok {
		storage.data.Store(userID, []string{noteToken})
		return nil
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
		return ErrFindUser
	}
	NotesIDs := rawNotesIDs.([]string)

	noteIndex, ok := findNote(NotesIDs, noteToken)
	if !ok {
		return ErrFindTokenInUsersNotes
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
