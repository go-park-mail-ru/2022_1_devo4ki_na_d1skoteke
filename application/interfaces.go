package application

import (
	"cotion/domain/entity"
)

type UserAppManager interface {
	SaveUser(registerUser entity.RegisterUserRequest) (entity.User, error)
	GetUser(email string) (entity.User, error)
}
