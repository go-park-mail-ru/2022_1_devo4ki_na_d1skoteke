package psql

import (
	"cotion/internal/api/domain/entity"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

const packageName = "psql"

type UserStorage struct {
	DB *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		DB: db,
	}
}

const querySaveUser = "INSERT INTO cotionuser(userid, username, email, password, avatar) VALUES ($1, $2, $3, $4, $5)"

func (store *UserStorage) Save(user entity.User) error {
	if _, err := store.DB.Exec(querySaveUser, user.UserID, user.Username, user.Email, user.Password, user.Avatar); err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Save",
		}).Error(err)
		return err
	}
	return nil
}

const queryGetUser = "SELECT userid, username, email, password, avatar FROM cotionuser WHERE userid = $1"

func (store *UserStorage) Get(userID string) (entity.User, error) {
	row := store.DB.QueryRow(queryGetUser, userID)
	user := entity.User{}
	if err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Avatar); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

const queryUpdateUser = "UPDATE cotionuser SET username = $1, password = $2, avatar = $3 where userid = $4"

func (store *UserStorage) Update(user entity.User) error {
	if _, err := store.DB.Exec(queryUpdateUser, user.Username, user.Password, user.Avatar, user.UserID); err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Update",
		}).Error(err)
		return err
	}
	return nil
}

const queryDeleteUser = "DELETE FROM cotionuser where userid = $1"

func (store *UserStorage) Delete(userID string) error {
	if _, err := store.DB.Exec(queryDeleteUser, userID); err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Delete",
		}).Error(err)
		return err
	}
	return nil
}
