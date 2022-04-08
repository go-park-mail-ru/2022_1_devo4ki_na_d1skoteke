package storage

import (
	"cotion/internal/domain/entity"
	"errors"
	"sync"
)

var ErrNoNoteInDB = errors.New("no note in DB with this token")

type NotesStorage struct {
	data sync.Map
}

func NewNotesStorage() *NotesStorage {
	store := &NotesStorage{
		data: sync.Map{},
	}
	store.data.Store("1", entity.Note{Name: "1st note", Body: "Hello everybody. This is Body of the 1st note)"})
	store.data.Store("2", entity.Note{Name: "2st note", Body: "Hello everybody. This is Body of the 2st note)"})
	store.data.Store("3", entity.Note{Name: "3st note", Body: "Hello everybody. This is Body of the 3st note)"})

	return store
}

func (store *NotesStorage) Find(token string) (entity.Note, error) {
	rawNote, ok := store.data.Load(token)
	if !ok {
		return entity.Note{}, ErrNoNoteInDB
	}
	note := rawNote.(entity.Note)
	return note, nil
}

func (store *NotesStorage) Save(token string, note entity.Note) error {
	_, ok := store.data.LoadOrStore(token, note)
	if !ok {
		return errors.New("there is note in DB with this token")
	}
	return nil
}

func (store *NotesStorage) Update(token string, note entity.Note) error {
	store.data.Store(token, note)
	return nil
}

func (store *NotesStorage) Delete(token string) error {
	store.data.Delete(token)
	return nil
}
