package main

import "database/sql"

//SessionView contains a group of workers to return to a client
type SessionView struct {
	Workers []ClientView `json:"workers"`
}

//Create creates a new session in the database
func (view *SessionView) Create(numWorkers int, clientID string, db *sql.DB) error {

	sessionID, err := GenUUID()
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
	row := db.QueryRow("COUNT (SELECT sessionID FROM session WHERE sessionID = ?)", id)
	err := row.Scan(&numWorkers)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT pubKey, address FROM worker WHERE clientID = ANY (SELECT workerID FROM session WHERE sessionID = ?)", id)
	if err != nil {
		return err
	}
	defer rows.Close()

	workers := make([]ClientView, numWorkers)
	view.Workers = workers

	i := 0
	for rows.Next() {
		err := rows.Scan(&view.Workers[i].PubKey, &view.Workers[i].Address)
		if err != nil {
			return err
		}
	}

	return nil
}
