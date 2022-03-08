package repository

import "cotion/domain/entity"

type SessionRepository interface {
	HasSession(SID string) (entity.Session, bool)
	NewSession(SID string, user entity.User) (entity.Session, error)
	DeleteSession(SID string)
}
