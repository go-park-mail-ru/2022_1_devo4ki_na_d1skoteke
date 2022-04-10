package storage

import (
	"cotion/internal/domain/entity"
	"cotion/internal/pkg/security"
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
	store.data.Store(string(manager.Hash("test2@mail.ru")), &entity.User{
		Username: "test2",
		Email:    "test2@mail.ru",
		Password: string(manager.Hash("Test1234!@#")),
	})
	store.data.Store(string(manager.Hash("nikita@mail.ru")), &entity.User{
		Username: "nikita",
		Email:    "nikita@mail.ru",
		Password: string(manager.Hash("Nikita1234!@#")),
	})
	return store
}

func (r *UserCacheStorage) Save(user entity.User) (entity.User, error) {
	r.data.LoadOrStore(string(r.securityManager.Hash(user.Email)), &user)
	return user, nil
}

func (r *UserCacheStorage) Get(email string) (entity.User, error) {
	rawUser, ok := r.data.Load(string(r.securityManager.Hash(email)))
	if ok {
		user := rawUser.(*entity.User)
		return *user, nil
	}
	return entity.User{}, errors.New("no user")
}

func (r *UserCacheStorage) Update(user entity.User) (entity.User, error) {
	r.data.Store(string(r.securityManager.Hash(user.Email)), user)
	return user, nil
}

func (r *UserCacheStorage) Delete(user entity.User) error {
	r.data.Delete(string(r.securityManager.Hash(user.Email)))
	return nil
}
