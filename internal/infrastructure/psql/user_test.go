package psql

import (
	"cotion/internal/domain/entity"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

var mockUser = entity.User{
	UserID:   "101",
	Username: "test",
	Email:    "test@mail.ru",
	Password: "Test1234!@#",
	Avatar:   "none",
}

func TestSaveUser(t *testing.T) {
	cases := map[string]struct {
		prepare  func(sqlmock.Sqlmock)
		expected func(error)
	}{
		"Success": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("INSERT INTO cotionuser").
					WithArgs(mockUser.UserID, mockUser.Username, mockUser.Email, mockUser.Password, mockUser.Avatar).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
		"Error": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("INSERT INTO cotionuser").
					WithArgs(mockUser.UserID, mockUser.Username, mockUser.Email, mockUser.Password, mockUser.Avatar).
					WillReturnError(fmt.Errorf("internal error"))
			},
			expected: func(actualErr error) {
				require.Equal(t, fmt.Errorf("internal error"), actualErr)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock)
			err := repo.Save(mockUser)
			tc.expected(err)
		})
		log.Println("SUCCESS")
	}
}

func TestGetUser(t *testing.T) {
	cases := map[string]struct {
		prepare  func(sqlmock.Sqlmock)
		expected func(entity.User, error)
	}{
		"Success": {
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"UserID", "Username", "Email", "Password", "Avatar"})
				rows = rows.AddRow(mockUser.UserID, mockUser.Username, mockUser.Email, mockUser.Password, mockUser.Avatar)
				mock.
					ExpectQuery("SELECT userid, username, email, password, avatar FROM cotionuser WHERE").
					WithArgs(mockUser.UserID).
					WillReturnRows(rows)
			},
			expected: func(actualNote entity.User, actualErr error) {
				require.Equal(t, nil, actualErr)
				require.Equal(t, mockUser, actualNote)
			},
		},
		"Error": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery("SELECT userid, username, email, password, avatar FROM cotionuser WHERE").
					WithArgs(mockUser.UserID).
					WillReturnError(fmt.Errorf("internal error"))
			},
			expected: func(actualNote entity.User, actualErr error) {
				require.Equal(t, fmt.Errorf("internal error"), actualErr)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock)
			user, err := repo.Get(mockUser.UserID)
			tc.expected(user, err)
		})
		log.Println("SUCCESS")
	}
}

func TestUpdateUser(t *testing.T) {
	cases := map[string]struct {
		prepare  func(sqlmock.Sqlmock)
		expected func(error)
	}{
		"Success": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("UPDATE cotionuser SET").
					WithArgs(mockUser.Username, mockUser.Password, mockUser.Avatar, mockUser.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
		"Error": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("UPDATE cotionuser SET").
					WithArgs(mockUser.Username, mockUser.Password, mockUser.Avatar, mockUser.UserID).
					WillReturnError(fmt.Errorf("internal error"))
			},
			expected: func(actualErr error) {
				require.Equal(t, fmt.Errorf("internal error"), actualErr)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock)
			err := repo.Update(mockUser)
			tc.expected(err)
		})
		log.Println("SUCCESS")
	}
}

func TestDeleteUser(t *testing.T) {
	cases := map[string]struct {
		prepare  func(sqlmock.Sqlmock)
		expected func(error)
	}{
		"Success": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("DELETE FROM cotionuser").
					WithArgs(mockUser.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: func(actualErr error) {
				require.Equal(t, nil, actualErr)
			},
		},
		"Error": {
			prepare: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec("DELETE FROM cotionuser").
					WithArgs(mockUser.UserID).
					WillReturnError(fmt.Errorf("internal error"))
			},
			expected: func(actualErr error) {
				require.Equal(t, fmt.Errorf("internal error"), actualErr)
			},
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserStorage(db)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.prepare(mock)
			err := repo.Delete(mockUser.UserID)
			tc.expected(err)
		})
		log.Println("SUCCESS")
	}
}
