package notes

import (
	"cotion/internal/domain/entity"
	"cotion/internal/infrastructure/storage"
	"cotion/internal/pkg/security"
	"errors"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestFindByToken(t *testing.T) {
	cases := map[string]struct {
		in       string
		expected func(entity.Note, error)
	}{
		"Note in DB": {
			in: "1",
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, entity.Note{
					Name: "1st note",
					Body: "Hello everybody. This is Body of the 1st note)",
				}, actualNote)
			},
		},
		"Note not in DB": {
			in: "0",
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, errors.New("no note in DB with this token"), actualErr)
				require.Equal(t, entity.Note{}, actualNote)
			},
		},
	}

	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	notesService := NewNotesApp(notesStorage, usersNotesStorage)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			result, err := notesService.FindByToken(tc.in)
			tc.expected(result, err)

		})
		log.Println("SUCCESS")
	}
}

func TestAllNotesByUserID(t *testing.T) {
	cases := map[string]struct {
		in       string
		expected func([]entity.Note, error)
	}{
		"User's notes in DB": {
			in: string(security.Hash("test@mail.ru")),
			expected: func(actualNote []entity.Note, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, []entity.Note{{
					Name: "1st note",
					Body: "Hello everybody. This is Body of the 1st note)",
				}, {
					Name: "3st note",
					Body: "Hello everybody. This is Body of the 3st note)",
				},
				}, actualNote)
			},
		},
		"No one user's note in DB": {
			in: string(security.Hash("test2@mail.ru")),
			expected: func(actualNote []entity.Note, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, []entity.Note{}, actualNote)
			},
		},
		"User has note token, but there isn't this note in DB": {
			in: string(security.Hash("test3@mail.ru")),
			expected: func(actualNote []entity.Note, actualErr error) {
				require.Equal(t, errors.New("cannot find note by token"), actualErr)
				require.Equal(t, []entity.Note{}, actualNote)
			},
		},
	}

	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	notesService := NewNotesApp(notesStorage, usersNotesStorage)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			result, err := notesService.AllNotesByUserID(tc.in)
			tc.expected(result, err)

		})
		log.Println("SUCCESS")
	}
}

func TestTokensByUserID(t *testing.T) {
	cases := map[string]struct {
		in       string
		expected func([]string, error)
	}{
		"There is user's tokens of notes": {
			in: string(security.Hash("test@mail.ru")),
			expected: func(actualNote []string, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, []string{"1", "3"}, actualNote)
			},
		},
		"There isn't user's tokens of notes": {
			in: string(security.Hash("test2@mail.ru")),
			expected: func(actualNote []string, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, []string{}, actualNote)
			},
		},
	}

	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	notesService := NewNotesApp(notesStorage, usersNotesStorage)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			result, err := notesService.TokensByUserID(tc.in)
			tc.expected(result, err)

		})
		log.Println("SUCCESS")
	}
}
