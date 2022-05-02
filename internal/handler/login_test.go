package handler

import (
	"bytes"
	mock_application "cotion/internal/application/mocks"
	"cotion/internal/domain/entity"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler_Login(t *testing.T) {
	cases := map[string]struct {
		inData       string
		mockBehavior func(s *mock_application.MockAuthAppManager)
		expected     func(w *httptest.ResponseRecorder)
	}{
		"Success": {
			inData: `{ 
				"email": "test@mail.ru",
				"password": "Test1234"
				}`,
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
				s.EXPECT().Login("test@mail.ru", "Test1234").Return(&http.Cookie{}, nil)
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 200, w.Code)
			},
		},
		"Empty email or password": {
			inData: `{ 
				"email": "",
				"password": "Test1234"
				}`,
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
				//s.EXPECT().Login("test@mail.ru", "Test1234").Return(&http.Cookie{}, nil)
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 400, w.Code)
				require.Equal(t, ErrNoLoginData.Error()+"\n", w.Body.String())
			},
		},
		"Wrong email or password": {
			inData: `{ 
				"email": "test@mail.ru",
				"password": "Test1234"
				}`,
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
				s.EXPECT().Login("test@mail.ru", "Test1234").Return(&http.Cookie{}, errors.New("wrong login data"))
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 400, w.Code)
				require.Equal(t, "wrong login data\n", w.Body.String())
			},
		},
		"Data is not json": {
			inData: ` 
				"email": "test@mail.ru",
				"password": "Test1234"
				`,
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 400, w.Code)
				require.Equal(t, ErrDecode.Error()+"\n", w.Body.String())
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_application.NewMockAuthAppManager(c)
			tc.mockBehavior(auth)

			loginHandler := NewLoginHandler(auth)

			router := mux.NewRouter()
			router.HandleFunc("/api/v1/users/login", loginHandler.Login).Methods("POST")

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/users/login",
				bytes.NewBufferString(tc.inData))

			router.ServeHTTP(w, req)

			tc.expected(w)
		})
		log.Println("SUCCESS")
	}
}

func TestLoginHandler_Logout(t *testing.T) {
	var mockCookie = &http.Cookie{Name: sessionCookie, Value: "aaaaaaa"}
	cases := map[string]struct {
		inCookie     *http.Cookie
		mockBehavior func(s *mock_application.MockAuthAppManager)
		expected     func(w *httptest.ResponseRecorder)
	}{
		"Success": {
			inCookie: mockCookie,
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
				s.EXPECT().
					Logout(mockCookie).
					Return(mockCookie, nil)
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 200, w.Code)
			},
		},
		"No cookie": {
			inCookie: &http.Cookie{},
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 401, w.Code)
			},
		},
		"No session": {
			inCookie: mockCookie,
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
				s.EXPECT().
					Logout(mockCookie).
					Return(nil, errors.New("no find session"))
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 401, w.Code)
				require.Equal(t, "no find session\n", w.Body.String())
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_application.NewMockAuthAppManager(c)
			tc.mockBehavior(auth)

			loginHandler := NewLoginHandler(auth)

			router := mux.NewRouter()
			router.HandleFunc("/api/v1/users/logout", loginHandler.Logout).Methods("GET")

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/users/logout", nil)
			req.AddCookie(tc.inCookie)

			router.ServeHTTP(w, req)

			tc.expected(w)
		})
		log.Println("SUCCESS")
	}
}

func TestLoginHandler_Auth(t *testing.T) {
	var mockCookie = &http.Cookie{Name: sessionCookie, Value: "aaaaaaa"}
	cases := map[string]struct {
		inCookie     *http.Cookie
		mockBehavior func(s *mock_application.MockAuthAppManager)
		expected     func(w *httptest.ResponseRecorder)
	}{
		"Success": {
			inCookie: mockCookie,
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
				s.EXPECT().
					Auth(mockCookie).
					Return(entity.User{}, true)
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 200, w.Code)
			},
		},
		"No cookie": {
			inCookie: &http.Cookie{},
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 401, w.Code)
			},
		},
		"No session": {
			inCookie: mockCookie,
			mockBehavior: func(s *mock_application.MockAuthAppManager) {
				s.EXPECT().
					Auth(mockCookie).
					Return(entity.User{}, false)
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 401, w.Code)
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_application.NewMockAuthAppManager(c)
			tc.mockBehavior(auth)

			loginHandler := NewLoginHandler(auth)

			router := mux.NewRouter()
			router.HandleFunc("/api/v1/users/auth", loginHandler.Auth).Methods("GET")

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/users/auth", nil)
			req.AddCookie(tc.inCookie)

			router.ServeHTTP(w, req)

			tc.expected(w)
		})
		log.Println("SUCCESS")
	}
}
