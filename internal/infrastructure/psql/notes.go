package psql

import (
	"cotion/internal/domain/entity"
	"database/sql"
	"errors"
)

var ErrNoNoteInDB = errors.New("no note in DB with this token")

type NotesStorage struct {
	DB *sql.DB
}

func NewNotesStorage(db *sql.DB) *NotesStorage {
	return &NotesStorage{
		DB: db,
	}
}

func (store *NotesStorage) Find(token string) (entity.Note, error) {
	row := store.DB.QueryRow("SELECT name, body FROM note WHERE NoteID = $1", token)
	note := &entity.Note{}
	if err := row.Scan(&note.Name, &note.Body); err != nil {
		return entity.Note{}, err
	}
	return *note, nil
}

func (store *NotesStorage) Save(token string, note entity.Note) error {
	_, err := store.DB.Exec(
		"INSERT INTO note(noteID, name, body) VALUES ($1, $2, $3)",
		token,
		note.Name,
		note.Body)
	return err
}

func (store *NotesStorage) Update(token string, note entity.Note) error {
	_, err := store.DB.Exec(
		"UPDATE note SET name = $1, body = $2 WHERE noteID = $3",
		note.Name,
		note.Body,
		token)
	return err
}

func (store *NotesStorage) Delete(token string) error {
	_, err := store.DB.Exec(
		"DELETE FROM note where noteid = $1",
		token)
	return err
}
