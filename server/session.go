package main

import "database/sql"

//SessionView contains a group of workers to return to a client
type SessionView struct {
	Workers []ClientView `json:"workers"`
}

//Create creates a new session in the database
func (view *SessionView) Create(numWorkers int, clientID string, db *sql.DB) error {
	return nil
}

//Get retreives associated workers with a session in the database
func (view *SessionView) Get(id string, db *sql.DB) error {

	return nil
}
