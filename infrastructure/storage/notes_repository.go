package storage

import (
	"cotion/domain/entity"
	"errors"
	"sync"
)

type NotesStorage struct {
	data sync.Map
}

func NewNotesStorage() *NotesStorage {
	store := NotesStorage{
		data: sync.Map{},
	}
	store.data.Store("1", entity.Note{Name: "1st note", Body: "Hello everybody. This is Body of the 1st note)"})
	store.data.Store("2", entity.Note{Name: "2st note", Body: "Hello everybody. This is Body of the 2st note)"})
	store.data.Store("3", entity.Note{Name: "3st note", Body: "Hello everybody. This is Body of the 3st note)"})
	return &store
}

func (store *NotesStorage) GetByToken(token string) (entity.Note, error) {
	rawNote, ok := store.data.Load(token)
	if !ok {
		return entity.Note{}, errors.New("no note in DB with this token")
	}
	note := rawNote.(entity.Note)
	return note, nil
}
