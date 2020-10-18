package main

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

//ErrNoRowsAffected is used to indicate no change has been made to the database
var ErrNoRowsAffected = errors.New("no rows affected")

//GenUUID generates an UUIDv4 according to RFC 4122
func GenUUID() (string, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

//Connect is just a wrapper around go sql connect and ping
func Connect(dsn string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
