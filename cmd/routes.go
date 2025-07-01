package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/VsenseTechnologies/biometric_http_server/internals/handlers"
	middlewares "github.com/VsenseTechnologies/biometric_http_server/internals/middleware"
	"github.com/VsenseTechnologies/biometric_http_server/repository"
	"github.com/gorilla/mux"
)

func InitilizeHttpRouters(db *sql.DB) http.Handler {

	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)

	adminHandler := handlers.NewAdminHandler(repository.NewAdminRepo(db))
	userHandler := handlers.NewUserHandler(repository.NewUserRepo(db))
	biometricHandler := handlers.NewBiometricHandler(repository.NewBiometricRepo(db))
	studentHandler := handlers.NewStudentHandler(repository.NewStudentRepo(db))
	excelHandler := handlers.NewExcelController(repository.NewExcelRepository(db))

	router.HandleFunc("/root/create/admin", adminHandler.CreateAdminHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/admin/login", adminHandler.AdminLogin).Methods("POST", "OPTIONS")
	router.Handle("/admin/access/user/{user_id}", middlewares.AuthMiddleware(http.HandlerFunc(userHandler.GiveUserAccessHandler))).Methods("POST", "OPTIONS")
	router.Handle("/admin/access/user", middlewares.AuthMiddleware(http.HandlerFunc(userHandler.GiveUserAccessHandler))).Methods("POST", "OPTIONS")
	router.Handle("/admin/get/users", middlewares.AuthMiddleware(http.HandlerFunc(userHandler.GetAllUsersHandler))).Methods("GET", "OPTIONS")

	router.HandleFunc("/user/register", userHandler.CreateUserHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/login", userHandler.UserLoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/update/password", userHandler.UpdateNewPasswordHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/forgotpassword", userHandler.ForgotPasswordHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/validate/otp", userHandler.ValidateOtpHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/update/time", userHandler.UpdateTime).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/unit_ids/{user_id}", userHandler.GetBiometricDevicesForRegisterForm).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/get/student_unit_ids/{unit_id}", userHandler.GetStudentUnitIdsForRegisterForm).Methods("GET", "OPTIONS")

	router.HandleFunc("/user/create/biometric_device", biometricHandler.CreateBiometricDeviceHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/delete/biometric_device/{unit_id}", biometricHandler.DeleteBiometricDeviceHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/update/biometric_device/label", biometricHandler.UpdateBiometricLabelHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/get/biometric_device/{user_id}", biometricHandler.GetBiometricDevicesHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/user/create/student", studentHandler.CreateNewStudentHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/update/student", studentHandler.UpdateStudentDetailsHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/get/student/{unit_id}", studentHandler.GetStudentDetailsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/delete/student", studentHandler.DeleteStudentHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/student/logs/{student_id}", studentHandler.GetStudentLogsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/student/download/pdf", studentHandler.DownloadPdfHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/user/student/download/excel", excelHandler.DownloadExcel).Methods("POST", "OPTIONS")

	tempDir := "/tmp"
	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}

	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(tempDir)))).Methods("GET", "OPTIONS")

	return router
}
