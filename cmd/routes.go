package main

import (
	"database/sql"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/handlers"
	"github.com/VsenseTechnologies/biometric_http_server/repository"
	"github.com/gorilla/mux"
)

func InitilizeHttpRouters(db *sql.DB) http.Handler {
	router := mux.NewRouter() 
	authHandler := handlers.NewAuthHandler(repository.NewAuthRepo(db))
	router.HandleFunc("/{id}/register" , authHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/{id}/login" , authHandler.LoginHandler).Methods("POST")

	return router
}