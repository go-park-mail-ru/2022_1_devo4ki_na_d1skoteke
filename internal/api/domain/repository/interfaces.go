package repository

import (
	"cotion/internal/api/domain/entity"
	"github.com/minio/minio-go/v7"
)

type SessionRepository interface {
	HasSession(SID string) (entity.Session, bool)
	NewSession(SID string, user entity.User) (entity.Session, error)
	DeleteSession(SID string)
}

type UserRepository interface {
	Save(user entity.User) error
	Get(userID string) (entity.User, error)
	Update(user entity.User) error
	Delete(userID string) error
}

type UsersNotesRepository interface {
	AddLink(userID string, noteToken string) error
	DeleteLink(userID string, noteToken string) error
	CheckLink(userID string, noteToken string) bool
	AllNotesByUserID(hashedEmail string) ([]entity.Note, error)
}

type NotesRepository interface {
	Save(token string, note entity.Note) error
	Update(token string, note entity.Note) error
	Delete(token string) error
	Find(token string) (entity.Note, error)
}

type ImageRepository interface {
	UploadFile(image entity.ImageUnit) (string, error)
	DownloadFile(imageID string) (*minio.Object, error)
}
