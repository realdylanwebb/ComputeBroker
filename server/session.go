package main

import (
	"database/sql"
	"log"
)

//SessionView contains a group of workers to return to a client
type SessionView struct {
	Workers      []ClientView `json:"workers"`
	SessionToken string       `json:"token"`
}

//AllSessionsView contains a group of user sessions
type AllSessionsView struct {
	Sessions []SessionView `json:"sessions"`
}

//Create creates a new session in the database
func (view *SessionView) Create(numWorkers int, clientID string, db *sql.DB, keys *KeyChain) error {

	sessionID, err := GenUUID()
	if err != nil {
		return err
	}

	claims := new(PrivateClaims)
	claims.ID = sessionID
	view.SessionToken, err = keys.Sign(claims)
	if err != nil {
		return err
	}

	ids := make([]string, numWorkers)
	rows, err := db.Query("SELECT clientID FROM client WHERE jobsAvailable > 0 LIMIT ?", numWorkers)
	if err != nil {
		return err
	}

	i := 0
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ids[i])
		i++
		if err != nil {
			ids[i] = ""
			continue
		}
	}

	//Create many to many relation for a session
	for i = 0; i < numWorkers; i++ {
		result, err := db.Exec("INSERT INTO session (clientID, workerID, sessionID) VALUES (?, ?, ?)", clientID, ids[i], sessionID)
		if err != nil {
			return err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		log.Print(affected)

		if affected == 0 {
			return ErrNoRowsAffected
		}
	}

	//Retrieve the worker information
	err = view.Get(sessionID, db)
	if err != nil {
		return err
	}

	return nil
}

//Get retreives associated workers with a session in the database
func (view *SessionView) Get(id string, db *sql.DB) error {

	var numWorkers int64
	row := db.QueryRow("SELECT COUNT(sessionID) FROM session WHERE sessionID = ?", id)
	err := row.Scan(&numWorkers)
	if err != nil {
		return err
	}

	log.Print(numWorkers)

	rows, err := db.Query("SELECT pubKey, address FROM client WHERE clientID IN (SELECT workerID FROM session WHERE sessionID = ?)", id)
	if err != nil {
		return err
	}
	defer rows.Close()

	workers := make([]ClientView, numWorkers)
	view.Workers = workers

	i := 0
	for rows.Next() {
		err := rows.Scan(&view.Workers[i].PubKey, &view.Workers[i].Address)
		i++
		if err != nil {
			return err
		}
	}

	return nil
}

//GetUserSessions returns all the sessions associated with a user
func (view *AllSessionsView) GetUserSessions(userID string, db *sql.DB) error {

	var count int64
	row := db.QueryRow("SELECT COUNT(DISTINCT sessionID) FROM session WHERE clientID = ?", userID)
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	sessions := make([]SessionView, count)
	ids := make([]string, count)
	rows, err := db.Query("SELECT DISTINCT sessionID FROM session WHERE clientID = ?", userID)

	i := 0
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ids[i])
		i++
		if err != nil {
			return err
		}
	}

	for i = 0; int64(i) < count; i++ {
		err := sessions[i].Get(ids[i], db)
		if err != nil {
			return err
		}
	}

	view.Sessions = sessions
	return nil
}
