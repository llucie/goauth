package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	InitCache()

	// Initialize router
	router := mux.NewRouter()

	// Set up router endpoints
	router.Methods("PUT").Path("/signin").HandlerFunc(Signup)
	router.Methods("POST").Path("/signin").HandlerFunc(Signin)
	router.Methods("POST").Path("/refresh").HandlerFunc(Refresh)
	router.Methods("GET").Path("/welcome").HandlerFunc(Welcome)

	// Listen on port 8080
	log.Fatal(http.ListenAndServe(":4242", router))
}
