package storage

import (
	"cotion/domain/entity"
	"cotion/infrastructure/security"
	"errors"
	"sync"
)

type UsersNotesStorage struct {
	data  sync.Map
	notes *NotesStorage
}

var cannotFindNotesForUser = errors.New("cannot find notes")

func NewUsersNotesStorage(notesStorage *NotesStorage) *UsersNotesStorage {
	storage := &UsersNotesStorage{
		data:  sync.Map{},
		notes: notesStorage,
	}
	storage.data.Store(string(security.Hash("test@mail.ru")), []string{"1", "3"})
	storage.data.Store(string(security.Hash("nikita@mail.ru")), []string{"2"})
	return storage
}

func (storage *UsersNotesStorage) AllNotesByUserID(hashedEmail string) ([]entity.Note, error) {
	rawNotesIDs, ok := storage.data.Load(hashedEmail)
	if ok != true {
		return []entity.Note{}, cannotFindNotesForUser
	}

	notesIDs := rawNotesIDs.([]string)
	notes := make([]entity.Note, 0)

	for _, id := range notesIDs {
		note, err := storage.notes.FindByToken(id)
		if err != nil {
			return []entity.Note{}, cannotFindNotesForUser
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (storage *UsersNotesStorage) TokensByUserID(hashedEmail string) ([]string, error) {
	rawNotesIDs, ok := storage.data.Load(hashedEmail)
	if ok != true {
		return []string{}, cannotFindNotesForUser
	}
	return rawNotesIDs.([]string), nil
}
