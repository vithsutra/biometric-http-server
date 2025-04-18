package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/VsenseTechnologies/biometric_http_server/internals/handlers"
	"github.com/VsenseTechnologies/biometric_http_server/repository"
	"github.com/gorilla/mux"
)

func InitilizeHttpRouters(db *sql.DB) http.Handler {
	router := mux.NewRouter()

	adminHandler := handlers.NewAdminHandler(repository.NewAdminRepo(db))
	userHandler := handlers.NewUserHandler(repository.NewUserRepo(db))
	biometricHandler := handlers.NewBiometricHandler(repository.NewBiometricRepo(db))
	studentHandler := handlers.NewStudentHandler(repository.NewStudentRepo(db))
	excelHandler := handlers.NewExcelController(repository.NewExcelRepository(db))

	router.HandleFunc("/root/create/admin", adminHandler.CreateAdminHandler).Methods("POST")

	router.HandleFunc("/admin/login", adminHandler.AdminLogin).Methods("POST")
	router.HandleFunc("/admin/create/user", userHandler.CreateUserHandler).Methods("POST")
	router.HandleFunc("/admin/access/user", userHandler.GiveUserAccessHandler).Methods("POST")
	router.HandleFunc("/admin/get/users", userHandler.GetAllUsersHandler).Methods("GET")

	router.HandleFunc("/user/login", userHandler.UserLoginHandler).Methods("POST")
	router.HandleFunc("/user/update/password", userHandler.UpdateNewPasswordHandler).Methods("POST")
	router.HandleFunc("/user/forgotpassword", userHandler.ForgotPasswordHandler).Methods("POST")
	router.HandleFunc("/user/validate/otp", userHandler.ValidateOtpHandler).Methods("POST")
	router.HandleFunc("/user/update/time", userHandler.UpdateTime).Methods("POST")
	router.HandleFunc("/user/unit_ids/{user_id}", userHandler.GetBiometricDevicesForRegisterForm).Methods("GET")
	router.HandleFunc("/user/get/student_unit_ids/{unit_id}", userHandler.GetStudentUnitIdsForRegisterForm).Methods("GET")

	router.HandleFunc("/user/create/biometric_device", biometricHandler.CreateBiometricDeviceHandler).Methods("POST")
	router.HandleFunc("/user/delete/biometric_device/{unit_id}", biometricHandler.DeleteBiometricDeviceHandler).Methods("GET")
	router.HandleFunc("/user/update/biometric_device/label", biometricHandler.UpdateBiometricLabelHandler).Methods("POST")
	router.HandleFunc("/user/get/biometric_device/{user_id}", biometricHandler.GetBiometricDevicesHandler).Methods("GET")
	router.HandleFunc("/user/clear/biometric_device/data", biometricHandler.ClearBiometricDeviceDataHandler).Methods("POST")

	router.HandleFunc("/user/create/student", studentHandler.CreateNewStudentHandler).Methods("POST")
	router.HandleFunc("/user/update/student", studentHandler.UpdateStudentDetailsHandler).Methods("POST")
	router.HandleFunc("/user/get/student/{unit_id}", studentHandler.GetStudentDetailsHandler).Methods("GET")
	router.HandleFunc("/user/delete/student", studentHandler.DeleteStudentHandler).Methods("POST")
	router.HandleFunc("/user/student/logs/{student_id}", studentHandler.GetStudentLogsHandler).Methods("GET")
	router.HandleFunc("/user/student/download/pdf", studentHandler.DownloadPdfHandler).Methods("POST")

	router.HandleFunc("/user/student/download/excel", excelHandler.DownloadExcel).Methods("POST")

	tempDir := "/Users/bunny/Desktop/Vithsutra"

	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(tempDir)))).Methods("GET")
	return router
}
