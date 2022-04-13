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
		inUserID    string
		inNoteToken string
		expected    func(entity.Note, error)
	}{
		"Note in DB": {
			inUserID:    security.Hash("test@mail.ru"),
			inNoteToken: "1",
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, entity.Note{
					Name: "1st note",
					Body: "Hello everybody. This is Body of the 1st note)",
				}, actualNote)
			},
		},
		"Note not in DB or not access": {
			inUserID:    security.Hash("test@mail.ru"),
			inNoteToken: "0",
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, ErrNoteAccess, actualErr)
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
			result, err := notesService.GetNote(tc.inUserID, tc.inNoteToken)
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
			in: security.Hash("test2@mail.ru"),
			expected: func(actualNote []entity.Note, actualErr error) {
				require.Equal(t, storage.ErrFindNotesForUser, actualErr)
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

func TestSaveNote(t *testing.T) {
	cases := map[string]struct {
		inUserID      string
		inNoteRequest entity.NoteRequest
		expected      func(error)
	}{
		"Success": {
			inUserID: security.Hash("test@mail.ru"),
			inNoteRequest: entity.NoteRequest{
				Name: "Test Save",
				Body: "This is body for test save note",
			},
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
	}

	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	notesService := NewNotesApp(notesStorage, usersNotesStorage)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			err := notesService.SaveNote(tc.inUserID, tc.inNoteRequest)
			tc.expected(err)

		})
		log.Println("SUCCESS")
	}
}

func TestGetNote(t *testing.T) {
	cases := map[string]struct {
		inUserID    string
		inNoteToken string
		expected    func(entity.Note, error)
	}{
		"Success": {
			inUserID:    string(security.Hash("test@mail.ru")),
			inNoteToken: "1",
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, entity.Note{
					Name: "1st note",
					Body: "Hello everybody. This is Body of the 1st note)",
				}, actualNote)
			},
		},
		"ErrNoteAccess": {
			inUserID:    string(security.Hash("test@mail.ru")),
			inNoteToken: "2",
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, ErrNoteAccess, actualErr)
			},
		},
		"ErrNoNoteInDB": {
			inUserID:    string(security.Hash("test@mail.ru")),
			inNoteToken: "0",
			expected: func(actualNote entity.Note, actualErr error) {
				require.Equal(t, storage.ErrNoNoteInDB, actualErr)
			},
		},
	}

	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	notesService := NewNotesApp(notesStorage, usersNotesStorage)

	usersNotesStorage.AddLink(string(security.Hash("test@mail.ru")), "0")

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			note, err := notesService.GetNote(tc.inUserID, tc.inNoteToken)
			tc.expected(note, err)

		})
		log.Println("SUCCESS")
	}
}

func TestUpdateNote(t *testing.T) {
	cases := map[string]struct {
		inUserID      string
		inNoteToken   string
		inNoteRequest entity.NoteRequest
		expected      func(error)
	}{
		"Success": {
			inUserID:    string(security.Hash("test@mail.ru")),
			inNoteToken: "1",
			inNoteRequest: entity.NoteRequest{
				Name: "Updated 1st note",
				Body: "Hello everybody. This is Body of the updated 1st note)",
			},
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
		"ErrNoteAccess": {
			inUserID:    string(security.Hash("test@mail.ru")),
			inNoteToken: "2",
			inNoteRequest: entity.NoteRequest{
				Name: "Updated 1st note",
				Body: "Hello everybody. This is Body of the updated 1st note)",
			},
			expected: func(actualErr error) {
				require.Equal(t, ErrNoteAccess, actualErr)
			},
		},
	}

	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	notesService := NewNotesApp(notesStorage, usersNotesStorage)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			err := notesService.UpdateNote(tc.inUserID, tc.inNoteToken, tc.inNoteRequest)
			tc.expected(err)

		})
		log.Println("SUCCESS")
	}
}

func TestDeleteNote(t *testing.T) {
	cases := map[string]struct {
		inUserID    string
		inNoteToken string
		expected    func(error)
	}{
		"Success": {
			inUserID:    string(security.Hash("test@mail.ru")),
			inNoteToken: "1",
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
	}

	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	notesService := NewNotesApp(notesStorage, usersNotesStorage)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			err := notesService.DeleteNote(tc.inUserID, tc.inNoteToken)
			tc.expected(err)

		})
		log.Println("SUCCESS")
	}
}
