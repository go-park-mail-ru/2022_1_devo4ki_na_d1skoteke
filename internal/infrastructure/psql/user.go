package psql

import (
	"cotion/internal/domain/entity"
	"database/sql"
)

type UserStorage struct {
	DB *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		DB: db,
	}
}

const querySaveUser = "INSERT INTO cotionuser(userid, username, email, password) VALUES ($1, $2, $3, $4)"

func (store *UserStorage) Save(user entity.User) error {
	_, err := store.DB.Exec(querySaveUser, user.UserID, user.Username, user.Email, user.Password)
	return err
}

const queryGetUser = "SELECT userid, username, email, password from cotionuser where userid = $1"

func (store *UserStorage) Get(userID string) (entity.User, error) {
	row := store.DB.QueryRow(queryGetUser, userID)
	user := entity.User{}
	if err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

const queryUpdateUser = "UPDATE cotionuser SET username = $1, password = $2 where userid = $3"

func (store *UserStorage) Update(user entity.User) error {
	_, err := store.DB.Exec(queryUpdateUser, user.Username, user.Password, user.UserID)
	return err
}

const queryDeleteUser = "DELETE FROM cotionuser where userid = $1"

func (store *UserStorage) Delete(userID string) error {
	_, err := store.DB.Exec(queryDeleteUser, userID)
	return err
}
