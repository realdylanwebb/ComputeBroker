package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//SessionRequest is used to unmarshall a session request body
type SessionRequest struct {
	Workers int `json:"workers"`
}

//Login returns a API key to a registered client
func (serv *BrokerServer) Login(w http.ResponseWriter, r *http.Request) {

	view := new(ClientView)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&view)
	if err != nil {
		log.Print(err)
		respondErr(w, 400, "Bad request body.")
		return
	}
}

//Register creates a new client
func (serv *BrokerServer) Register(w http.ResponseWriter, r *http.Request) {

	view := new(ClientView)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&view)
	if err != nil {
		log.Print(err)
		respondErr(w, 400, "Bad request body.")
		return
	}

}

//Signal changes the amount of jobs that a client is ready to recieve
func (serv *BrokerServer) Signal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	claims, err := serv.keys.Validate(r.Header.Get("Authorization"))
	if err != nil {
		log.Print(err)
		respondErr(w, 401, "Invalid API key.")
		return
	}

	available, err := strconv.Atoi(vars["available"])

	view := new(ClientView)
	view.JobsAvailable = int64(available)
	view.ClientID = claims.ID

	view.Signal(serv.DB)
}

//ReqSession creates a new session and returns the associated worker information and
//a session key
func (serv *BrokerServer) ReqSession(w http.ResponseWriter, r *http.Request) {

	body := new(SessionRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Print(err)
		respondErr(w, 400, "Bad request body.")
		return
	}

	claims, err := serv.keys.Validate(r.Header.Get("Authorization"))
	if err != nil {
		log.Print(err)
		respondErr(w, 401, "Invalid API key.")
		return
	}

	view := new(SessionView)
	err = view.Create(body.Workers, claims.ID, serv.DB)
	if err != nil {
		log.Print(err)
		respondErr(w, 500, "Internal server error.")
		return
	}

	respondJSON(w, 201, view)
	return
}

//GetSession gets the worker information associated with a session id
func (serv *BrokerServer) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	claims, err := serv.keys.Validate(vars["id"])
	if err != nil {
		log.Print(err)
		respondErr(w, 401, "Invalid API key.")
		return
	}

	view := new(SessionView)
	err = view.Get(claims.ID, serv.DB)
	if err != nil {
		log.Print(err)
		respondErr(w, 500, "Internal server error.")
		return
	}

	respondJSON(w, 200, view)
	return
}

////////////////////
//RESPONSE HELPERS//
////////////////////
func respondErr(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
	return
}

func respondOK(w http.ResponseWriter, msg string) {
	respondJSON(w, 200, map[string]string{"success": msg})
	return
}

func respondJSON(w http.ResponseWriter, status int, res interface{}) {
	response, err := json.Marshal(res)
	if err != nil {
		log.Print(err)
		respondErr(w, http.StatusInternalServerError, "Internal server error.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
	return
}
