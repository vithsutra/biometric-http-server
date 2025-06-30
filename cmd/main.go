package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	serverMode := os.Getenv("SERVER_MODE")
	serverMode = "dev"

	if serverMode == "dev" {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("missing the .env file", err)
		}
		log.Println("running in development mode..")
		log.Println(".env file loaded successfully")
		return
	}

	if serverMode == "prod" {
		log.Println("running in the production mode")
		log.Println(".env file loading skipped")
		return
	}

	log.Fatalln("please set the SERVER_MODE to 'dev' for development use or to 'prod' for production use")
}

func main() {
	// connecting to database and checking the status of connection
	db := NewDatabase()
	db.CheckStatus()
	defer db.Close()
	// Starting the server
	Start(db.db)
}
