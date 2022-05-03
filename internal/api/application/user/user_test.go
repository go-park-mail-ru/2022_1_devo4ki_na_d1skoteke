package user

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/infrastructure/storage"
	"cotion/internal/pkg/security"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	username              = "username"
	email                 = "email@mail.ru"
	badEmail              = "emailru"
	securePassword        = "6Va}4HG%/h+fHUYD"
	anotherSecurePassword = "6Va}4HG%/h+fHUYD23sd4!@3213"
	unsecurePassword      = "password"
)

func TestUserApplication(t *testing.T) {
	cases := map[string]struct {
		prepare  func() (*UserService, entity.UserRequest)
		process  func(userService *UserService, requestUser entity.UserRequest) (*UserService, entity.UserRequest, error)
		expected func(*UserService, entity.UserRequest, error)
	}{
		"save and get user": {
			prepare: func() (*UserService, entity.UserRequest) {
				securityManager := security.NewSimpleSecurityManager()
				userStorage := storage.NewUserCacheStorage(securityManager)
				var imageStorage repository.ImageRepository
				requestUser := entity.UserRequest{
					Username:        username,
					Email:           email,
					Password:        securePassword,
					ConfirmPassword: securePassword,
				}
				return NewUserService(userStorage, imageStorage, securityManager), requestUser
			},
			process: func(userService *UserService, requestUser entity.UserRequest) (*UserService, entity.UserRequest, error) {
				err := userService.Save(requestUser)
				return userService, requestUser, err
			},
			expected: func(userService *UserService, requestUser entity.UserRequest, err error) {
				require.NoError(t, err)

				user, err := userService.Get(security.Hash(requestUser.Email))
				require.NoError(t, err)

				require.Equal(t, requestUser.Username, user.Username)
				require.Equal(t, requestUser.Email, user.Email)
				require.Equal(t, string(userService.securityManager.Hash(requestUser.Password)), user.Password)
				require.Equal(t, string(userService.securityManager.Hash(requestUser.ConfirmPassword)), user.Password)
			},
		},
		"save registered user": {
			prepare: func() (*UserService, entity.UserRequest) {
				securityManager := security.NewSimpleSecurityManager()
				userStorage := storage.NewUserCacheStorage(securityManager)
				var imageStorage repository.ImageRepository
				requestUser := entity.UserRequest{
					Username:        username,
					Email:           email,
					Password:        securePassword,
					ConfirmPassword: securePassword,
				}
				return NewUserService(userStorage, imageStorage, securityManager), requestUser
			},
			process: func(userService *UserService, requestUser entity.UserRequest) (*UserService, entity.UserRequest, error) {
				err := userService.Save(requestUser)
				require.NoError(t, err)
				err = userService.Save(requestUser)

				return userService, requestUser, err
			},
			expected: func(userService *UserService, requestUser entity.UserRequest, err error) {
				require.ErrorIs(t, err, ErrUserAlreadyRegistered)
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			userService, requestUser := tc.prepare()
			userService, user, err := tc.process(userService, requestUser)
			tc.expected(userService, user, err)
		})
	}
}
