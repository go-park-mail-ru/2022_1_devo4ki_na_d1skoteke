package application

import (
	"cotion/infrastructure/security"
	"cotion/infrastructure/storage"
	"errors"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	cases := map[string]struct {
		inParam1 string
		inParam2 string
		expected func(*http.Cookie, error)
	}{
		"Successful login": {
			inParam1: "test@mail.ru",
			inParam2: "Test1234!@#",
			expected: func(actualCookie *http.Cookie, actualErr error) {
				require.Equal(t, nil, actualErr)
				expectedCookie := &http.Cookie{
					Name: sessionCookie,
					Path: pathSessionCookie,
				}
				require.Equal(t, expectedCookie.Name, actualCookie.Name)
				require.Equal(t, expectedCookie.Path, actualCookie.Path)
			},
		},
		"Incorrect email": {
			inParam1: "test0@mail.ru",
			inParam2: "Test1234!@#",
			expected: func(actualCookie *http.Cookie, actualErr error) {
				require.Equal(t, errors.New("no user"), actualErr)
				require.Equal(t, &http.Cookie{}, actualCookie)
			},
		},
		"Incorrect password": {
			inParam1: "test@mail.ru",
			inParam2: "#Test1234!@#",
			expected: func(actualCookie *http.Cookie, actualErr error) {
				require.Equal(t, errors.New("wrong password"), actualErr)
				require.Equal(t, &http.Cookie{}, actualCookie)
			},
		},
	}

	securityManager := security.NewSimpleSecurityManager()
	sessionStorage := storage.NewSessionStorage()
	userStorage := storage.NewUserCacheStorage(securityManager)

	userService := NewUserService(userStorage, securityManager)
	authService := NewAuthApp(sessionStorage, userService, securityManager)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			result, err := authService.Login(tc.inParam1, tc.inParam2)
			tc.expected(result, err)
		})
		log.Println("SUCCESS")
	}
}

func TestLoginLogout(t *testing.T) {
	cases := map[string]struct {
		inParam1 string
		inParam2 string
		expected func(*http.Cookie, error)
	}{
		"Successful": {
			inParam1: "test@mail.ru",
			inParam2: "Test1234!@#",
			expected: func(actualCookie *http.Cookie, actualErr error) {
				require.Equal(t, nil, actualErr)
				expectedCookie := &http.Cookie{
					Name: sessionCookie,
					Path: pathSessionCookie,
				}
				require.Equal(t, expectedCookie.Name, actualCookie.Name)
				require.Equal(t, expectedCookie.Path, actualCookie.Path)
			},
		},
		"No session": {
			inParam1: "test0@mail.ru",
			inParam2: "Test1234!@#",
			expected: func(actualCookie *http.Cookie, actualErr error) {
				require.Equal(t, errors.New("no session"), actualErr)
				require.Equal(t, &http.Cookie{}, actualCookie)
			},
		},
	}

	securityManager := security.NewSimpleSecurityManager()
	sessionStorage := storage.NewSessionStorage()
	userStorage := storage.NewUserCacheStorage(securityManager)

	userService := NewUserService(userStorage, securityManager)
	authService := NewAuthApp(sessionStorage, userService, securityManager)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			sessionCookie, _ := authService.Login(tc.inParam1, tc.inParam2)
			result, err := authService.Logout(*sessionCookie)
			tc.expected(result, err)
		})
		log.Println("SUCCESS")
	}
}
