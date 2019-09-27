package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

// Store the redis connection as a package level variable
var cache redis.Conn

var users = map[string]string{}

// Credentials is a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// InitCache is a function
func InitCache() {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	// Assign the connection to the package level `cache` variable
	cache = conn
}

func createSession(w http.ResponseWriter, username string) {
	// Create a new random session token
	sessionToken, err := uuid.NewUUID()
	if err != nil {
		// Fail to generate an UUID, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 5 minutes
	_, err = cache.Do("SETEX", sessionToken, "300", username)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 5 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: time.Now().Add(5 * time.Minute),
	})
}

func getSession(w http.ResponseWriter, r *http.Request) (string, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return "", errors.New("Username not recognized")
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return "", errors.New("Bad request")
	}
	sessionToken := c.Value

	log.Print("Session token is ", sessionToken)

	// We then get the name of the user from our cache, where we set the session token
	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		return "", errors.New("Internal error")
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("Unauthorized")
	}

	return fmt.Sprintf("%s", response), nil
}
