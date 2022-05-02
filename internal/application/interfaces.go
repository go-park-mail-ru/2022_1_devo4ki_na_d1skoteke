package application

import (
	"cotion/internal/domain/entity"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"net/http"
)

type NotesAppManager interface {
	//FindByToken(token string) (entity.Note, error)
	AllNotesByUserID(userID string) (entity.ShortNotes, error)
	SaveNote(userID string, noteRequest entity.NoteRequest) error
	GetNote(userID string, noteToken string) (entity.Note, error)
	UpdateNote(userID string, noteToken string, noteRequest entity.NoteRequest) error
	DeleteNote(userID string, noteToken string) error
}

type UserAppManager interface {
	Save(registerUser entity.UserRequest) error
	Get(userID string) (entity.User, error)
	Update(curUser entity.User, user entity.UserRequest) error
	Delete(userID string) error
	UploadAvatar(src multipart.File, hdr *multipart.FileHeader, user entity.User) error
	DownloadAvatar(user entity.User) (*minio.Object, error)
}

type AuthAppManager interface {
	Login(login string, password string) (*http.Cookie, error)
	Logout(sessionCookie *http.Cookie) (*http.Cookie, error)
	Auth(sessionCookie *http.Cookie) (entity.User, bool)
}
