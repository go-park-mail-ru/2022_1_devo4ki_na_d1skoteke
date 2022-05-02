package handler

import (
	mock_application "cotion/internal/application/mocks"
	"cotion/internal/domain/entity"
	"cotion/internal/pkg/security"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotesHandler_ReceiveSingleNote(t *testing.T) {
	var mockUser = entity.User{
		UserID:   "c04532ca4e12438bcd37d2ae1676d3f5a27241062095eaccdbf0102b78d2a948",
		Email:    "test@mail.ru",
		Password: "Test1234",
	}
	var mockToken = "1"
	var mockNote = entity.Note{
		Name: "Test note",
		Body: "Body of test note",
	}
	cases := map[string]struct {
		inToken      string
		mockBehavior func(notesServ *mock_application.MockNotesAppManager)
		expected     func(w *httptest.ResponseRecorder)
	}{
		"Success": {
			inToken: mockToken,
			mockBehavior: func(notesServ *mock_application.MockNotesAppManager) {
				notesServ.EXPECT().
					GetNote(mockUser.UserID, mockToken).
					Return(mockNote, nil)
			},
			expected: func(w *httptest.ResponseRecorder) {
				require.Equal(t, 200, w.Code)
				//require.Equal(t, "some", w.Body.String())
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			notesServ := mock_application.NewMockNotesAppManager(c)
			authServ := mock_application.NewMockAuthAppManager(c)
			securityServ := security.NewSimpleSecurityManager()

			tc.mockBehavior(notesServ)

			notesHandler := NewNotesHandler(notesServ, authServ, securityServ)

			router := mux.NewRouter()
			router.HandleFunc("/api/v1/note/{note-token:[0-9]+}", notesHandler.ReceiveSingleNote).Methods("GET")

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/note/"+tc.inToken, nil)
			req = req.WithContext(context.WithValue(req.Context(), "user", mockUser))

			router.ServeHTTP(w, req)

			tc.expected(w)
		})
		log.Println("SUCCESS")
	}
}
