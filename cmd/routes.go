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
	adminHandler := handlers.NewAdminHandler(repository.NewAdminRepo(db))
	router.HandleFunc("/admin/getusers" , adminHandler.FetchAllUsersHandler).Methods("GET")
	router.HandleFunc("/admin/giveaccess" , adminHandler.GiveUserAccessHandler).Methods("POST")
	biometricHandler := handlers.NewBiometricHandler(repository.NewBiometricRepo(db))
	router.HandleFunc("/admin/devices/{userid}" , biometricHandler.FetchAllBiometricsHandler).Methods("GET")
	router.HandleFunc("/admin/device/delete/{unitid}" , biometricHandler.DeleteBiometricMachineHandler).Methods("GET")
	router.HandleFunc("/admin/adddevice" , biometricHandler.NewBiometricDeviceHandler).Methods("POST")
	return router
}