package s3

import (
	"cotion/internal/domain/entity"
	"cotion/internal/pkg/generator"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
)

const (
	ENV_MINIO_URL      = "minio_url"
	ENV_MINIO_USER     = "minio_user"
	ENV_MINIO_PASSWORD = "minio_pass"
	BucketName         = "avatars"
	packageName        = "s3"
)

var ErrNoUrl = errors.New("There isn't minio url in *.env file")
var ErrNoUser = errors.New("There isn't minio user in *.env file")
var ErrNoPass = errors.New("There isn't minio password in *.env file")

type MinioProvider struct {
	client *minio.Client
}

func NewMinioProvider() (*MinioProvider, error) {
	var urlEnv, userEnv, passwordEnv string
	if urlEnv = os.Getenv(ENV_MINIO_URL); urlEnv == "" {
		return &MinioProvider{}, ErrNoUrl
	}
	if userEnv = os.Getenv(ENV_MINIO_USER); userEnv == "" {
		return &MinioProvider{}, ErrNoUser
	}
	if passwordEnv = os.Getenv(ENV_MINIO_PASSWORD); passwordEnv == "" {
		return &MinioProvider{}, ErrNoPass
	}

	client, err := minio.New(urlEnv, &minio.Options{
		Creds:  credentials.NewStaticV4(userEnv, passwordEnv, ""),
		Secure: false,
	})
	if err != nil {
		return &MinioProvider{}, err
	}

	ctx := context.Background()

	err = client.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, BucketName)
		if errBucketExists == nil && exists {
			log.Info("We already own ", BucketName)
		} else {
			return &MinioProvider{}, err
		}
	} else {
		log.Info("Successfully created ", BucketName)
	}

	return &MinioProvider{
		client,
	}, nil
}

//
//func (m *MinioProvider) Connect() error {
//	var err error
//
//	return err
//}

func (m *MinioProvider) UploadFile(unit entity.ImageUnit) (string, error) {
	imageName := generator.RandSID(16) + ".png"

	_, err := m.client.PutObject(
		context.Background(),
		BucketName,
		imageName,
		unit.Payload,
		unit.PayloadSize,
		minio.PutObjectOptions{ContentType: "image/png"},
	)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "UploadFile",
		}).Error(err)
	}

	return imageName, err
}

func (m *MinioProvider) DownloadFile(imageName string) (*minio.Object, error) {
	reader, err := m.client.GetObject(
		context.Background(),
		BucketName,
		imageName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "DownloadFile",
		}).Error(err)
		return nil, err
	}

	return reader, nil
}
