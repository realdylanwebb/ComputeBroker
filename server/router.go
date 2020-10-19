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

type TokenResponse struct {
	Token string `json:"token"`
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

	res := new(TokenResponse)

	res.Token, err = view.Token(serv.DB, serv.keys)
	if err != nil {
		log.Print(err)
		respondErr(w, 401, "Invalid credentials.")
		return
	}
	respondJSON(w, 200, res)
	return
}

//Register creates a new client
func (serv *BrokerServer) Register(w http.ResponseWriter, r *http.Request) {
	log.Print("HERE")
	view := new(ClientView)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&view)
	if err != nil {
		log.Print(err)
		respondErr(w, 400, "Bad request body.")
		return
	}

	err = view.Create(serv.DB)
	if err != nil {
		if err.Error()[:7] == "UNIQUE " {
			respondErr(w, 400, "Client already exists with that email")
			return
		}
		log.Print(err)
		respondErr(w, 500, "Internal server error.")
		return
	}
	respondJSON(w, 201, view)
	return
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

	err = view.Signal(serv.DB)
	if err != nil {
		log.Print(err)
		respondErr(w, 500, "Internal server error.")
		return
	}
	respondOK(w, "Signaled readyness.")
	return
}

//GetUser returns user information to a valid user
func (serv *BrokerServer) GetUser(w http.ResponseWriter, r *http.Request) {
	claims, err := serv.keys.Validate(r.Header.Get("Authorization"))
	if err != nil {
		log.Print(err)
		respondErr(w, 401, "Invalid API key.")
		return
	}

	view := new(ClientView)
	view.ClientID = claims.ID

	err = view.Read(serv.DB)
	if err != nil {
		log.Print(err)
		respondErr(w, 500, "Internal server error.")
		return
	}
	respondJSON(w, 200, view)
	return
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
		log.Print(r.Header.Get("Authorization"))
		log.Print(err)
		respondErr(w, 401, "Invalid API key.")
		return
	}

	view := new(SessionView)
	err = view.Create(body.Workers, claims.ID, serv.DB, serv.keys)
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

	body := new(TokenResponse)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Print(err)
		respondErr(w, 400, "Bad request body.")
		return
	}

	claims, err := serv.keys.Validate("Bearer " + body.Token)
	if err != nil {
		log.Print(err)
		respondErr(w, 401, "Invalid API key.")
		return
	}

	_, err = serv.keys.Validate(r.Header.Get("Authorization"))
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

//GetUserSessions retrieves all the sessions associated with a worker
func (serv *BrokerServer) GetUserSessions(w http.ResponseWriter, r *http.Request) {

	claims, err := serv.keys.Validate(r.Header.Get("Authorization"))
	if err != nil {
		log.Print(err)
		respondErr(w, 401, "Invalid API key.")
		return
	}

	userSessions := new(AllSessionsView)
	err = userSessions.GetUserSessions(claims.ID, serv.DB)
	if err != nil {
		log.Print(err)
		respondErr(w, 500, "Internal server error.")
		return
	}

	respondJSON(w, 200, userSessions)
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
