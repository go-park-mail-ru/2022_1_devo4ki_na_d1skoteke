package user

import (
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/pkg/security"
	"errors"
)

const packageName = "app user"

type UserService struct {
	userRepository  repository.UserRepository
	securityManager security.Manager
}

var ErrUserAlreadyRegistered = errors.New("user already registered with this email")

func NewUserService(userRepository repository.UserRepository, securityManager security.Manager) *UserService {
	return &UserService{
		userRepository:  userRepository,
		securityManager: securityManager,
	}
}

func (u *UserService) Save(registerUser entity.UserRequest) error {
	user := entity.User{
		UserID:   u.securityManager.Hash(registerUser.Email),
		Username: registerUser.Username,
		Email:    registerUser.Email,
		Password: u.securityManager.Hash(registerUser.Password),
	}

	if _, err := u.userRepository.Get(user.UserID); err == nil {
		return ErrUserAlreadyRegistered
	}

	return u.userRepository.Save(user)
}

func (u *UserService) Get(userID string) (entity.User, error) {
	return u.userRepository.Get(userID)
}

func (u *UserService) Update(curUser entity.User, userRequest entity.UserRequest) error {
	user := entity.User{
		UserID:   curUser.UserID,
		Username: userRequest.Username,
		Email:    curUser.Email,
		Password: u.securityManager.Hash(userRequest.Password),
	}

	return u.userRepository.Update(user)
}

func (u *UserService) Delete(userID string) error {
	return u.userRepository.Delete(userID)
}
