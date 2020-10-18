package main

import (
	"database/sql"
	"errors"
)

//ClientView contains all potential fields for a user
type ClientView struct {
	ClientID      string
	Email         string `json:"email"`
	Password      string `json:"password"`
	PubKey        string `json:"pubKey"`
	Address       string `json:"address"`
	JobsAvailable int64  `json:"jobsAvailable"`
}

//Create inserts a client into the database
func (view *ClientView) Create(db *sql.DB) error {

	//generate client id
	id, err := GenUUID()
	if err != nil {
		return err
	}
	//hash password
	hashed := HashPass(view.Password)

	result, err := db.Exec("INSERT INTO client (clientID, email, password, pubKey, address, jobsAvailable) VALUES (?, ?, ?, ?, ?, ?)",
		id, view.Email, hashed, view.PubKey, view.Address, 0)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

//Token returns a client's API key
func (view *ClientView) Token(db *sql.DB, keys *KeyChain) (string, error) {
	var hashed string
	var clientID string
	row := db.QueryRow("SELECT clientID, password FROM client WHERE email = ?", view.Email)
	err := row.Scan(&clientID, &hashed)
	if err != nil {
		return "", err
	}

	if VerifyPass(view.Password, hashed) {
		claims := new(PrivateClaims)
		claims.ID = clientID
		return keys.Sign(claims)
	}

	return "", errors.New("invalid login")
}

//Signal is used to update the amount of jobs the client can accept
func (view *ClientView) Signal(db *sql.DB) error {

	result, err := db.Exec("UPDATE client SET jobsAvailable = ? WHERE clientID = ?", view.JobsAvailable, view.ClientID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}
