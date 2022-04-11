package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

const (
	host    = "postgres"
	port    = 5432
	dbname  = "cotion"
	userKey = "user"
	passKey = "pass"
)

var NoConnectData = errors.New("There isn't connect data in *.env file")

func ConnectConfig() (string, error) {
	user := os.Getenv(userKey)
	password := os.Getenv(passKey)
	if user == "" || password == "" {
		return "", NoConnectData
	}

	config := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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
