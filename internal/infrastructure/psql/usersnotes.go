package psql

import (
	"cotion/internal/domain/entity"
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
)

var ErrNoteAccess = errors.New("No link in usersnotes table")

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
	log.WithFields(log.Fields{
		"package":  packageName,
		"function": "AddLink",
		"userID":   userID,
		"noteID":   noteToken,
	}).Error(err)
	return err
}

const queryDeleteLink = "DELETE FROM usersnotes WHERE userid = $1 AND noteid = $2"

func (store *UsersNotesStorage) DeleteLink(userID string, noteToken string) error {
	_, err := store.DB.Exec(queryDeleteLink, userID, noteToken)
	log.WithFields(log.Fields{
		"package":  packageName,
		"function": "DeleteLink",
	}).Error(err)
	return err
}

const queryCheckLink = "SELECT 'exist' FROM usersnotes WHERE userid = $1 AND noteid = $2"

func (store *UsersNotesStorage) CheckLink(userID string, noteToken string) bool {
	row := store.DB.QueryRow(queryCheckLink, userID, noteToken)
	if row.Scan() == sql.ErrNoRows {
		log.WithFields(log.Fields{
			"package":   packageName,
			"function":  "CheckLink",
			"userID":    userID,
			"noteToken": noteToken,
		}).Warning(ErrNoteAccess)
		return false
	}
	return true
}

const queryFindNotes = "SELECT name, body, note.noteid FROM usersnotes JOIN note ON usersnotes.noteid = note.noteid WHERE userid = $1"

func (store *UsersNotesStorage) AllNotesByUserID(userID string) (entity.ShortNotes, error) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "AllNotesByUserID",
	})

	rows, err := store.DB.Query(queryFindNotes, userID)
	if err != nil {
		logger.Error(err)
		return entity.ShortNotes{}, err
	}
	defer rows.Close()

	notes := entity.ShortNotes{}
	for rows.Next() {
		var note entity.ShortNote
		if err := rows.Scan(&note.Name, &note.Body, &note.Token); err != nil {
			logger.Error(err)
			return entity.ShortNotes{}, err
		}
		notes.ShortNote = append(notes.ShortNote, note)
	}

	if err := rows.Err(); err != nil {
		logger.Error(err)
		return entity.ShortNotes{}, err
	}

	return notes, nil
}
