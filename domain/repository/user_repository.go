package repository

import "cotion/domain/entity"

type UserRepository interface {
	SaveUser(user entity.User) (entity.User, error)
	GetUser(email string) (entity.User, error)
}
