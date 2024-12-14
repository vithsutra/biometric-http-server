package main

import (
	"log"

	"github.com/joho/godotenv"
)


func main() {
	// Loading the environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Unable to load .env file: %v",err)
	}
	// connecting to database and checking the status of connection
	db := NewDatabase()
	db.CheckStatus()
	defer db.Close()
	// Starting the server
	Start(db.db)
}