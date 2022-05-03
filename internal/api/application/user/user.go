package user

import (
	"cotion/internal/api/domain/entity"
	"cotion/internal/api/domain/repository"
	"cotion/internal/pkg/security"
	"errors"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
)

const packageName = "app user"
const emptyAvatar = "none"

type UserService struct {
	userRepository  repository.UserRepository
	imageRepository repository.ImageRepository
	securityManager security.Manager
}

var ErrUserAlreadyRegistered = errors.New("user already registered with this email")
var ErrUserNoHasAvatar = errors.New("user hasn't avatar")

func NewUserService(userRepo repository.UserRepository, imageRepo repository.ImageRepository, securityManager security.Manager) *UserService {
	return &UserService{
		userRepository:  userRepo,
		imageRepository: imageRepo,
		securityManager: securityManager,
	}
}

func (u *UserService) Save(registerUser entity.UserRequest) error {
	user := entity.User{
		UserID:   u.securityManager.Hash(registerUser.Email),
		Username: registerUser.Username,
		Email:    registerUser.Email,
		Password: u.securityManager.Hash(registerUser.Password),
		Avatar:   emptyAvatar,
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

func (u *UserService) UploadAvatar(src multipart.File, hdr *multipart.FileHeader, user entity.User) error {
	object := entity.ImageUnit{
		Payload:     src,
		PayloadSize: hdr.Size,
	}
	defer src.Close()

	imageName, err := u.imageRepository.UploadFile(object)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "UploadAvatar",
		}).Error(err)
		return err
	}

	user.Avatar = imageName
	return u.userRepository.Update(user)
}

func (u *UserService) DownloadAvatar(user entity.User) (*minio.Object, error) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "DownloadAvatar",
	})
	if user.Avatar == emptyAvatar {
		logger.Debug(ErrUserNoHasAvatar)
		return nil, ErrUserNoHasAvatar
	}

	img, err := u.imageRepository.DownloadFile(user.Avatar)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return img, nil
}
