package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

const (
	ENV_PG_HOST     = "host"
	ENV_PG_PORT     = "port"
	ENV_PG_DBNAME   = "dbname"
	ENV_PG_USER     = "user"
	ENV_PG_PASSWORD = "pass"
)

var ErrNoHost = errors.New("There isn't host in *.env file")
var ErrNoPort = errors.New("There isn't port in *.env file")
var ErrNoDBName = errors.New("There isn't database name in *.env file")
var ErrNoUser = errors.New("There isn't user in *.env file")
var ErrNoPass = errors.New("There isn't password in *.env file")

func ConnectConfig() (string, error) {
	var host, port, dbname, user, pass string
	if host = os.Getenv(ENV_PG_HOST); host == "" {
		return "", ErrNoHost
	}
	if port = os.Getenv(ENV_PG_PORT); port == "" {
		return "", ErrNoPort
	}
	if dbname = os.Getenv(ENV_PG_DBNAME); dbname == "" {
		return "", ErrNoDBName
	}
	if user = os.Getenv(ENV_PG_USER); user == "" {
		return "", ErrNoUser
	}
	if pass = os.Getenv(ENV_PG_PASSWORD); pass == "" {
		return "", ErrNoPass
	}

	config := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)
	return config, nil
}

func Connect() (*sql.DB, error) {
	dbConfig, err := ConnectConfig()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
