package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
)

func Start(db *sql.DB) {
	// Creating new server struct
	server := &http.Server{
		Addr: os.Getenv("PORT"),
		Handler: InitilizeHttpRouters(db),
	}

	// Initilizing database with tables
	query := database.NewQuery(db)
	if err := query.InitilizeDatabase(); err != nil {
		log.Fatalf("Unable to initilize database: %v" , err)
	}

	// Running the server
	log.Printf("Server is running at port %s" , os.Getenv("PORT"))
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Unable to start server: %v" , err)
	}
}