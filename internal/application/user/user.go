package user

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/pkg/security"
	"errors"
)

type UserService struct {
	userRepository  repository.UserRepository
	securityManager security.Manager
}

var UserAlreadyRegistered = errors.New("user already registered")

func NewUserService(userRepository repository.UserRepository, securityManager security.Manager) *UserService {
	return &UserService{
		userRepository:  userRepository,
		securityManager: securityManager,
	}
}

func (u *UserService) Save(registerUser entity.RegisterUserRequest) (entity.User, error) {
	if err := registerUser.Validate(); err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		Username: registerUser.Username,
		Email:    registerUser.Email,
		Password: registerUser.Password,
	}
	user.Password = string(u.securityManager.Hash(user.Password))

	if _, err := u.userRepository.Get(user.Email); err == nil {
		return user, UserAlreadyRegistered
	}

	savedUser, err := u.userRepository.Save(user)
	if err != nil {
		return savedUser, err
	}
	return savedUser, nil
}

func (u *UserService) Get(email string) (entity.User, error) {
	return u.userRepository.Get(email)
}

func (u *UserService) Update(user entity.User) (entity.User, error) {
	return u.userRepository.Update(user)
}

func (u *UserService) Delete(user entity.User) error {
	return u.userRepository.Delete(user)
}
