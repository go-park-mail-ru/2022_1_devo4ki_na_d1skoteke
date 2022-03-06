package application

import (
	"cotion/domain/entity"
	"cotion/domain/repository"
	"cotion/infrastructure/security"
	"errors"
)

type UserService struct {
	userRepository  repository.UserRepository
	securityManager security.Manager
}

func NewUserService(userRepository repository.UserRepository, securityManager security.Manager) *UserService {
	return &UserService{
		userRepository:  userRepository,
		securityManager: securityManager,
	}
}

func (u *UserService) SaveUser(registerUser entity.RegisterUserRequest) (entity.User, error) {
	err := registerUser.Validate()
	if err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		Username: registerUser.Username,
		Email:    registerUser.Email,
		Password: registerUser.Password,
	}
	user.Password = string(u.securityManager.Hash(user.Password))

	_, err = u.userRepository.GetUser(user.Email)
	if err == nil {
		return user, errors.New("user already registered")
	}

	savedUser, err := u.userRepository.SaveUser(user)
	if err != nil {
		return savedUser, err
	}
	return savedUser, nil
}

func (u *UserService) GetUser(email string) (entity.User, error) {
	return u.userRepository.GetUser(email)
}
