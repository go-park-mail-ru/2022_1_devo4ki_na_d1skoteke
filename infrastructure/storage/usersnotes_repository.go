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
	storage.data.Store(string(security.Hash("email@vk.team")), []string{"1", "2", "3"})
	return storage
}

func (storage *UsersNotesStorage) GetAllNotesByUserID(hashedEmail string) ([]entity.Note, error) {
	rawNotesIDs, ok := storage.data.Load(hashedEmail)
	if ok != true {
		return []entity.Note{}, cannotFindNotesForUser
	}

	notesIDs := rawNotesIDs.([]string)
	notes := make([]entity.Note, 1)

	for _, id := range notesIDs {
		note, err := storage.notes.GetByToken(id)
		if err != nil {
			return []entity.Note{}, cannotFindNotesForUser
		}
		notes = append(notes, note)
	}

	return notes, nil
}
