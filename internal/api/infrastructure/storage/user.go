package storage

import (
	"cotion/internal/api/domain/entity"
	"cotion/internal/pkg/security"
	"errors"
	"sync"
)

var ErrNoUserInDB = errors.New("no user in database with this userID")
var ErrUserAlreadyRegistered = errors.New("user with this userID already has in database")

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

func (r *UserCacheStorage) Save(user entity.User) error {
	_, loaded := r.data.LoadOrStore(user.UserID, &user)
	if loaded {
		return ErrUserAlreadyRegistered
	}
	return nil
}

func (r *UserCacheStorage) Get(userID string) (entity.User, error) {
	rawUser, ok := r.data.Load(userID)
	if !ok {
		return entity.User{}, ErrNoUserInDB
	}
	user := rawUser.(*entity.User)
	return *user, nil
}

func (r *UserCacheStorage) Update(user entity.User) error {
	r.data.Store(user.UserID, &user)
	return nil
}

func (r *UserCacheStorage) Delete(userID string) error {
	if _, loaded := r.data.LoadAndDelete(userID); !loaded {
		return ErrNoUserInDB
	}
	return nil
}
