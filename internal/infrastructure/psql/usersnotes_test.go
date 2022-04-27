package psql

import (
	"cotion/internal/domain/entity"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

var mockNoteToken = "adjfkjanfkakdfjjk"
var mockNote = entity.Note{
	Name: "testNoteName",
	Body: "testNoteBody",
}

func TestCheckLink(t *testing.T) {
	cases := map[string]struct {
		prepare  func(sqlmock.Sqlmock)
		expected func(bool)
	}{
		"Ok": {
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exist"})
				rows = rows.AddRow("exist")
				mock.
					ExpectQuery("SELECT 'exist' FROM usersnotes WHERE").
					WithArgs(mockUser.UserID, mockNoteToken).
					WillReturnRows(rows)
			},
			expected: func(actualResult bool) {
				require.Equal(t, true, actualResult)
			},
		},
		"False": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery("SELECT 'exist' FROM usersnotes WHERE").
					WithArgs(mockUser.UserID, mockNoteToken).
					WillReturnError(sql.ErrNoRows)
			},
			expected: func(actualResult bool) {
				require.Equal(t, false, actualResult)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUsersNotesStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock)
			ok := repo.CheckLink(mockUser.UserID, mockNoteToken)
			tc.expected(ok)
		})
		log.Println("SUCCESS")
	}
}

func TestAddLink(t *testing.T) {
	cases := map[string]struct {
		prepare  func(sqlmock.Sqlmock)
		expected func(error)
	}{
		"Success": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("INSERT INTO usersnotes").
					WithArgs(mockUser.UserID, mockNoteToken).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: func(actualError error) {
				require.Equal(t, nil, actualError)
			},
		},
		"Error": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("INSERT INTO usersnotes").
					WithArgs(mockUser.UserID, mockNoteToken).
					WillReturnError(fmt.Errorf("internal error"))
			},
			expected: func(actualError error) {
				require.Equal(t, fmt.Errorf("internal error"), actualError)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUsersNotesStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock)
			err := repo.AddLink(mockUser.UserID, mockNoteToken)
			tc.expected(err)
		})
		log.Println("SUCCESS")
	}
}

func TestDeleteLink(t *testing.T) {
	cases := map[string]struct {
		prepare  func(sqlmock.Sqlmock)
		expected func(error)
	}{
		"Success": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("DELETE FROM usersnotes WHERE").
					WithArgs(mockUser.UserID, mockNoteToken).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: func(actualError error) {
				require.Equal(t, nil, actualError)
			},
		},
		"False": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("DELETE FROM usersnotes WHERE").
					WithArgs(mockUser.UserID, mockNoteToken).
					WillReturnError(fmt.Errorf("internal error"))
			},
			expected: func(actualError error) {
				require.Equal(t, fmt.Errorf("internal error"), actualError)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUsersNotesStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock)
			err := repo.DeleteLink(mockUser.UserID, mockNoteToken)
			tc.expected(err)
		})
		log.Println("SUCCESS")
	}
}

func TestAllNotesByUserID(t *testing.T) {
	cases := map[string]struct {
		prepare  func(sqlmock.Sqlmock)
		expected func([]entity.Note, error)
	}{
		"Success": {
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"Name", "Body"})
				rows = rows.AddRow(mockNote.Name, mockNote.Body)
				mock.
					ExpectQuery("SELECT name, body FROM usersnotes").
					WithArgs(mockUser.UserID).
					WillReturnRows(rows)
			},
			expected: func(actualResult []entity.Note, actualError error) {
				require.Equal(t, nil, actualError)
				require.Equal(t, actualResult, []entity.Note{mockNote})
			},
		},
		"Error": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery("SELECT name, body FROM usersnotes JOIN note").
					WithArgs(mockUser.UserID).
					WillReturnError(fmt.Errorf("internal error"))
			},
			expected: func(actualResult []entity.Note, actualError error) {
				require.Equal(t, fmt.Errorf("internal error"), actualError)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUsersNotesStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock)
			result, err := repo.AllNotesByUserID(mockUser.UserID)
			tc.expected(result, err)
		})
		log.Println("SUCCESS")
	}
}
