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
	studentHandler := handlers.NewStudentHandler(repository.NewStudentRepo(db))
	router.HandleFunc("/users/getstudents/{unitid}" , studentHandler.FetchStudentDetails).Methods("GET")
	router.HandleFunc("/users/getlogs/{studentid}" , studentHandler.FetchStudentLogHistoryHandler).Methods("GET")
	router.HandleFunc("/users/delete/{unitid}/{studentid}/{studentunitid}" , studentHandler.DeleteStudentHandler).Methods("GET")
	router.HandleFunc("/users/update" , studentHandler.UpdateStudentHandler).Methods("POST")
	router.HandleFunc("/users/generatepdf" , studentHandler.GenerateStudentAttendenceReportHandler).Methods("POST")
	router.HandleFunc("/users/newstudent" , studentHandler.NewStudentHandler).Methods("POST")
	excelHandler := handlers.NewExcelHandler(repository.NewExcelRepo(db))
	router.HandleFunc("/user/excel" , excelHandler.GenerateExcelReportHandler).Methods("POST")
	return router
}