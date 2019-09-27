package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Signup is a function
func Signup(w http.ResponseWriter, r *http.Request) {
	log.Print("Received AddUser request")
	var creds Credentials

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Add new user
	log.Print("Add user with Username = ", creds.Username)
	users[creds.Username] = creds.Password

	w.WriteHeader(http.StatusCreated)
}

// Signin is a function
func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	createSession(w, creds.Username)
}

// Refresh is a function
func Refresh(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user's session
	username := getSession(w, r)

	// Now, create a new session token for the current user
	createSession(w, username)
}

// Welcome is a function
func Welcome(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user's session
	username := getSession(w, r)

	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", username)))
}
