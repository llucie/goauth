package main

import "github.com/gomodule/redigo/redis"

// Store the redis connection as a package level variable
var cache redis.Conn

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
