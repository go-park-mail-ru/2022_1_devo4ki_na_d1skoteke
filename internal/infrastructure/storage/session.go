package storage

import (
	"cotion/internal/domain/entity"
	"errors"
	"sync"
)

type SessionStorage struct {
	data sync.Map
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{}
}

func (s *SessionStorage) NewSession(SID string, user entity.User) (entity.Session, error) {
	session, loaded := s.data.LoadOrStore(SID, entity.Session{
		SID:       SID,
		UserEmail: user.Email,
	})
	if loaded {
		return entity.Session{}, errors.New("already has session with this SID")
	}
	return session.(entity.Session), nil
}

func (s *SessionStorage) HasSession(SID string) (entity.Session, bool) {
	session, ok := s.data.Load(SID)
	if !ok {
		return entity.Session{}, false
	}
	return session.(entity.Session), true
}

func (s *SessionStorage) DeleteSession(SID string) {
	s.data.LoadAndDelete(SID)
}
