package psql

import (
	"cotion/internal/domain/entity"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestFindNotes(t *testing.T) {
	const (
		noteName = "name"
		noteBody = "body of the note"
	)

	cases := map[string]struct {
		inNoteToken string
		prepare     func(sqlmock.Sqlmock, string)
		expected    func(entity.Note, error)
	}{
		"Success": {
			inNoteToken: "1",
			prepare: func(mock sqlmock.Sqlmock, noteToken string) {
				rows := sqlmock.NewRows([]string{"name", "body"})
				rows = rows.AddRow(noteName, noteBody)
				mock.
					ExpectQuery("SELECT name, body FROM note WHERE").
					WithArgs(noteToken).
					WillReturnRows(rows)
			},
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, entity.Note{
					Name: noteName,
					Body: noteBody,
				}, actualNote)
			},
		},
		"No note in DB": {
			inNoteToken: "1",
			prepare: func(mock sqlmock.Sqlmock, noteToken string) {
			},
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, ErrNoNoteInDB, actualErr)

			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewNotesStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock, tc.inNoteToken)
			item, err := repo.Find(tc.inNoteToken)
			tc.expected(item, err)
		})
		log.Println("SUCCESS")
	}
}

func TestSaveNote(t *testing.T) {
	const (
		noteName = "name"
		noteBody = "body of the note"
	)

	cases := map[string]struct {
		inNoteToken string
		prepare     func(sqlmock.Sqlmock, string)
		expected    func(error)
	}{
		"Success": {
			inNoteToken: "1",
			prepare: func(mock sqlmock.Sqlmock, noteToken string) {
				mock.
					ExpectExec("INSERT INTO note").
					WithArgs(noteToken, noteName, noteBody).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
		"Err: already has note with this token": {
			inNoteToken: "1",
			prepare: func(mock sqlmock.Sqlmock, noteToken string) {
				mock.
					ExpectExec("INSERT INTO note").
					WithArgs(noteToken, noteName, noteBody).
					WillReturnError(fmt.Errorf("already has note with this token"))
			},
			expected: func(actualErr error) {
				require.Equal(t, fmt.Errorf("already has note with this token"), actualErr)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewNotesStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock, tc.inNoteToken)
			err := repo.Save(tc.inNoteToken, entity.Note{
				Name: noteName,
				Body: noteBody,
			})
			tc.expected(err)
		})
		log.Println("SUCCESS")
	}
}

func TestUpdateNote(t *testing.T) {
	const (
		noteName = "updated name"
		noteBody = "updated body of the note"
	)

	cases := map[string]struct {
		inNoteToken string
		prepare     func(sqlmock.Sqlmock, string)
		expected    func(error)
	}{
		"Success": {
			inNoteToken: "1",
			prepare: func(mock sqlmock.Sqlmock, noteToken string) {
				mock.
					ExpectExec("UPDATE note").
					WithArgs(noteName, noteBody, noteToken).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
		"No note in DB": {
			inNoteToken: "1",
			prepare: func(mock sqlmock.Sqlmock, noteToken string) {
				mock.
					ExpectExec("UPDATE note").
					WithArgs(noteName, noteBody, noteToken).
					WillReturnError(fmt.Errorf("no note with this token"))
			},
			expected: func(actualErr error) {
				require.Equal(t, fmt.Errorf("no note with this token"), actualErr)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewNotesStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock, tc.inNoteToken)
			err := repo.Update(tc.inNoteToken, entity.Note{
				Name: noteName,
				Body: noteBody,
			})
			tc.expected(err)
		})
		log.Println("SUCCESS")
	}
}

func TestDeleteNote(t *testing.T) {
	cases := map[string]struct {
		inNoteToken string
		prepare     func(sqlmock.Sqlmock, string)
		expected    func(error)
	}{
		"Success": {
			inNoteToken: "1",
			prepare: func(mock sqlmock.Sqlmock, noteToken string) {
				mock.
					ExpectExec("DELETE FROM note").
					WithArgs(noteToken).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
		"No note in DB": {
			inNoteToken: "1",
			prepare: func(mock sqlmock.Sqlmock, noteToken string) {
				mock.
					ExpectExec("DELETE FROM note").
					WithArgs(noteToken).
					WillReturnError(fmt.Errorf("no note with this token"))
			},
			expected: func(actualErr error) {
				require.Equal(t, fmt.Errorf("no note with this token"), actualErr)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewNotesStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock, tc.inNoteToken)
			err := repo.Delete(tc.inNoteToken)
			tc.expected(err)
		})
		log.Println("SUCCESS")
	}
}
