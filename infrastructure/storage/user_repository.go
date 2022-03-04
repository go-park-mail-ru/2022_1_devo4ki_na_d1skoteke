package storage

import (
	"cotion/domain/entity"
	"cotion/infrastructure/security"
	"errors"
	"sync"
)

type UserCacheStorage struct {
	data            sync.Map
	securityManager security.Manager
}

func NewUserCacheStorage(manager security.Manager) *UserCacheStorage {
	return &UserCacheStorage{
		data:            sync.Map{},
		securityManager: manager,
	}
}

func (r *UserCacheStorage) SaveUser(user entity.User) (entity.User, error) {
	r.data.LoadOrStore(r.securityManager.Hash(user.Email), &user)
	return user, nil
}

func (r *UserCacheStorage) GetUser(email string) (entity.User, error) {
	rawUser, ok := r.data.Load(r.securityManager.Hash(email))
	if ok {
		user := rawUser.(entity.User)
		return user, nil
	}
	return entity.User{}, errors.New("no user")
}
