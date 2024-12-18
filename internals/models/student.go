package models

import "net/http"

type Student struct {
	StudentId       string `json:"student_id,omitempty"`
	StudentUnitId   string `json:"student_unit_id,omitempty"`
	StudentName     string `json:"student_name,omitempty"`
	Department      string `json:"department,omitempty"`
	UnitId          string `json:"unit_id,omitempty"`
	FingerprintData string `json:"fingerprint_data,omitempty"`
	Date            string `json:"date,omitempty"`
	LoginTime       string `json:"login,omitempty"`
	LogoutTime      string `json:"logout,omitempty"`
}

type StudentInterface interface {
	NewStudent(*http.Request) error
	DeleteStudent(*http.Request) error
	UpdateStudent(*http.Request) error
	FetchStudentDetails(*http.Request) ([]Student, error)
	FetchStudentLogHistory(*http.Request) ([]Student, error)
    GenerateStudentAttendenceReport(*http.Request)
}
