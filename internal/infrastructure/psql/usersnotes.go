package psql

import (
	"cotion/internal/domain/entity"
	"database/sql"
	"errors"
)

var ErrFindNotesForUser = errors.New("cannot find notes")
var ErrFindNoteByToken = errors.New("cannot find note by token")
var ErrFindUser = errors.New("can't find user")
var ErrFindTokenInUsersNotes = errors.New("can't find token in user's notes")

type UsersNotesStorage struct {
	DB *sql.DB
}

func NewUsersNotesStorage(db *sql.DB) *UsersNotesStorage {
	return &UsersNotesStorage{
		DB: db,
	}
}

const queryAddLink = "INSERT INTO usersnotes(userid, noteid) VALUES ($1, $2)"

func (store *UsersNotesStorage) AddLink(userID string, noteToken string) error {
	_, err := store.DB.Exec(queryAddLink, userID, noteToken)
	return err
}

const queryDeleteLink = "DELETE FROM usersnotes WHERE userid = $1 AND noteid = $2"

func (store *UsersNotesStorage) DeleteLink(userID string, noteToken string) error {
	_, err := store.DB.Exec(queryDeleteLink, userID, noteToken)
	return err
}

const queryCheckLink = "SELECT 'exist' FROM usersnotes WHERE userid = $1 AND noteid = $2"

func (store *UsersNotesStorage) CheckLink(userID string, noteToken string) bool {
	row := store.DB.QueryRow(queryCheckLink, userID, noteToken)
	if row.Err() != nil {
		return false
	}
	return true
}

const queryFindNotes = "SELECT name, body FROM usersnotes JOIN note ON usersnotes.noteid = note.noteid WHERE userid = $1"

func (store *UsersNotesStorage) AllNotesByUserID(userID string) ([]entity.Note, error) {
	rows, err := store.DB.Query(queryFindNotes, userID)
	if err != nil {
		return []entity.Note{}, err
	}
	defer rows.Close()

	notes := make([]entity.Note, 0)
	for rows.Next() {
		var note entity.Note
		if err := rows.Scan(&note.Name, &note.Body); err != nil {
			return []entity.Note{}, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return []entity.Note{}, err
	}

	return notes, nil
}
