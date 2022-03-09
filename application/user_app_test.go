package application

import (
	"cotion/domain/entity"
	"cotion/infrastructure/security"
	"cotion/infrastructure/storage"
	"github.com/stretchr/testify/require"
	"testing"
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
		prepare  func() (*UserService, entity.RegisterUserRequest)
		process  func(userService *UserService, requestUser entity.RegisterUserRequest) (*UserService, entity.RegisterUserRequest, error)
		expected func(*UserService, entity.RegisterUserRequest, error)
	}{
		"save and get user": {
			prepare: func() (*UserService, entity.RegisterUserRequest) {
				securityManager := security.NewSimpleSecurityManager()
				userStorage := storage.NewUserCacheStorage(securityManager)
				requestUser := entity.RegisterUserRequest{
					Username:        username,
					Email:           email,
					Password:        securePassword,
					ConfirmPassword: securePassword,
				}
				return NewUserService(userStorage, securityManager), requestUser
			},
			process: func(userService *UserService, requestUser entity.RegisterUserRequest) (*UserService, entity.RegisterUserRequest, error) {
				_, err := userService.SaveUser(requestUser)
				return userService, requestUser, err
			},
			expected: func(userService *UserService, requestUser entity.RegisterUserRequest, err error) {
				require.NoError(t, err)

				user, err := userService.GetUser(requestUser.Email)
				require.NoError(t, err)

				require.Equal(t, requestUser.Username, user.Username)
				require.Equal(t, requestUser.Email, user.Email)
				require.Equal(t, string(userService.securityManager.Hash(requestUser.Password)), user.Password)
				require.Equal(t, string(userService.securityManager.Hash(requestUser.ConfirmPassword)), user.Password)
			},
		},
		"weak password": {
			prepare: func() (*UserService, entity.RegisterUserRequest) {
				securityManager := security.NewSimpleSecurityManager()
				userStorage := storage.NewUserCacheStorage(securityManager)
				requestUser := entity.RegisterUserRequest{
					Username:        username,
					Email:           email,
					Password:        unsecurePassword,
					ConfirmPassword: unsecurePassword,
				}
				return NewUserService(userStorage, securityManager), requestUser
			},
			process: func(userService *UserService, requestUser entity.RegisterUserRequest) (*UserService, entity.RegisterUserRequest, error) {
				_, err := userService.SaveUser(requestUser)
				return userService, requestUser, err
			},
			expected: func(userService *UserService, requestUser entity.RegisterUserRequest, err error) {
				require.Error(t, err)
			},
		},
		"bad email": {
			prepare: func() (*UserService, entity.RegisterUserRequest) {
				securityManager := security.NewSimpleSecurityManager()
				userStorage := storage.NewUserCacheStorage(securityManager)
				requestUser := entity.RegisterUserRequest{
					Username:        username,
					Email:           badEmail,
					Password:        unsecurePassword,
					ConfirmPassword: unsecurePassword,
				}
				return NewUserService(userStorage, securityManager), requestUser
			},
			process: func(userService *UserService, requestUser entity.RegisterUserRequest) (*UserService, entity.RegisterUserRequest, error) {
				_, err := userService.SaveUser(requestUser)
				return userService, requestUser, err
			},
			expected: func(userService *UserService, requestUser entity.RegisterUserRequest, err error) {
				require.Error(t, err)
			},
		},
		"different passwords": {
			prepare: func() (*UserService, entity.RegisterUserRequest) {
				securityManager := security.NewSimpleSecurityManager()
				userStorage := storage.NewUserCacheStorage(securityManager)
				requestUser := entity.RegisterUserRequest{
					Username:        username,
					Email:           email,
					Password:        unsecurePassword,
					ConfirmPassword: anotherSecurePassword,
				}
				return NewUserService(userStorage, securityManager), requestUser
			},
			process: func(userService *UserService, requestUser entity.RegisterUserRequest) (*UserService, entity.RegisterUserRequest, error) {
				_, err := userService.SaveUser(requestUser)
				return userService, requestUser, err
			},
			expected: func(userService *UserService, requestUser entity.RegisterUserRequest, err error) {
				require.ErrorIs(t, err, entity.DifferentPasswords)
			},
		},
		"save registered user": {
			prepare: func() (*UserService, entity.RegisterUserRequest) {
				securityManager := security.NewSimpleSecurityManager()
				userStorage := storage.NewUserCacheStorage(securityManager)
				requestUser := entity.RegisterUserRequest{
					Username:        username,
					Email:           email,
					Password:        securePassword,
					ConfirmPassword: securePassword,
				}
				return NewUserService(userStorage, securityManager), requestUser
			},
			process: func(userService *UserService, requestUser entity.RegisterUserRequest) (*UserService, entity.RegisterUserRequest, error) {
				_, err := userService.SaveUser(requestUser)
				require.NoError(t, err)
				_, err = userService.SaveUser(requestUser)

				return userService, requestUser, err
			},
			expected: func(userService *UserService, requestUser entity.RegisterUserRequest, err error) {
				require.ErrorIs(t, err, UserAlreadyRegistered)
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
