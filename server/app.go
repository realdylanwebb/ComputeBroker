package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

//BrokerServer contains the database connection and the router
type BrokerServer struct {
	DB     *sql.DB
	Router *mux.Router
	keys   *KeyChain
}

//Init creates a databse connection and initializes the router and keychain
func (serv *BrokerServer) Init() {

	dbName := os.Getenv("DBNAME")
	dbRef, err := Connect(dbName)
	if err != nil {
		panic(err)
	}

	serv.DB = dbRef
	serv.Router = mux.NewRouter()

	var keys KeyChain
	dur, err := time.ParseDuration(os.Getenv("TOKENTTL"))
	if err != nil {
		panic(err)
	}

	keys.Init(os.Getenv("AUTHENCKEY"), os.Getenv("AUTHSIGKEY"), dur)

	//Resource: User
	serv.Router.HandleFunc("/client", serv.Register).Methods("POST")
	serv.Router.HandleFunc("/client/signal/{available}", serv.Signal).Methods("POST")
	serv.Router.HandleFunc("/login", serv.Login).Methods("POST")

	//Resource: Session
	serv.Router.HandleFunc("/session", serv.ReqSession).Methods("POST")
	serv.Router.HandleFunc("/session/{key}", serv.GetSession).Methods("GET")

}

//Run starts the server on the specified address
func (serv *BrokerServer) Run() {
	var httpConfig http.Server
	httpConfig.Handler = serv.Router
	httpConfig.Addr = os.Getenv("SERVICEHOST")
	log.Fatal(httpConfig.ListenAndServe())
}

func main() {
	var serv BrokerServer
	serv.Init()
	serv.Run()
}
