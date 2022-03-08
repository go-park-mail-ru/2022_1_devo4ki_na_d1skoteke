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
	store := &UserCacheStorage{
		data:            sync.Map{},
		securityManager: manager,
	}
	store.data.Store(string(manager.Hash("test@mail.ru")), &entity.User{
		Username: "test",
		Email:    "test@mail.ru",
		Password: string(manager.Hash("Test1234!@#")),
	})
	store.data.Store(string(manager.Hash("nikita@mail.ru")), &entity.User{
		Username: "nikita",
		Email:    "nikita@mail.ru",
		Password: string(manager.Hash("Nikita1234!@#")),
	})
	return store
}

func (r *UserCacheStorage) SaveUser(user entity.User) (entity.User, error) {
	r.data.LoadOrStore(string(r.securityManager.Hash(user.Email)), &user)
	return user, nil
}

func (r *UserCacheStorage) GetUser(email string) (entity.User, error) {
	rawUser, ok := r.data.Load(string(r.securityManager.Hash(email)))
	if ok {
		user := rawUser.(*entity.User)
		return *user, nil
	}
	return entity.User{}, errors.New("no user")
}
