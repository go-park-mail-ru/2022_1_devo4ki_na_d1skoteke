package psql

import (
	"cotion/internal/domain/entity"
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
)

const packageName = "psql"

var ErrNoNoteInDB = errors.New("no note in DB with this token")

type NotesStorage struct {
	DB *sql.DB
}

func NewNotesStorage(db *sql.DB) *NotesStorage {
	return &NotesStorage{
		DB: db,
	}
}

const queryFindNote = "SELECT name, body FROM note WHERE NoteID = $1"

func (store *NotesStorage) Find(token string) (entity.Note, error) {
	row := store.DB.QueryRow(queryFindNote, token)
	note := entity.Note{}
	if err := row.Scan(&note.Name, &note.Body); err != nil {
		log.WithFields(log.Fields{
			"package":   packageName,
			"function":  "Find",
			"noteToken": token,
		}).Warning(err)
		return entity.Note{}, ErrNoNoteInDB
	}
	return note, nil
}

const querySaveNote = "INSERT INTO note(noteID, name, body) VALUES ($1, $2, $3)"

func (store *NotesStorage) Save(token string, note entity.Note) error {
	if _, err := store.DB.Exec(querySaveNote, token, note.Name, note.Body); err != nil {
		log.WithFields(log.Fields{
			"package":   packageName,
			"function":  "Save",
			"note":      note,
			"noteToken": token,
		}).Error(err)
		return err
	}
	return nil
}

const queryUpdateNote = "UPDATE note SET name = $1, body = $2 WHERE noteID = $3"

func (store *NotesStorage) Update(token string, note entity.Note) error {
	if _, err := store.DB.Exec(queryUpdateNote, note.Name, note.Body, token); err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Update",
		}).Error(err)
		return err
	}
	return nil
}

const queryDeleteNote = "DELETE FROM note where noteid = $1"

func (store *NotesStorage) Delete(token string) error {
	if _, err := store.DB.Exec(queryDeleteNote, token); err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Delete",
		}).Error(err)
		return err
	}
	return nil
}
