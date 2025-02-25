package models

import "net/http"

type CreateStudentRequest struct {
	StudentUnitId   string `json:"student_unit_id" validate:"required"`
	StudentName     string `json:"student_name" validate:"required"`
	StudentUsn      string `json:"student_usn" validate:"required"`
	Department      string `json:"department" validate:"required"`
	UnitId          string `json:"unit_id" validate:"required"`
	FingerprintData string `json:"fingerprint_data" validate:"required"`
}

type UpdateStudentRequest struct {
	UnitId      string `json:"unit_id" validate:"required"`
	StudentId   string `json:"student_id" validate:"required"`
	StudentName string `json:"student_name" validate:"required"`
	StudentUsn  string `json:"student_usn" validate:"required"`
	Department  string `json:"department" validate:"required"`
}

type DeleteStudentRequest struct {
	UnitId        string `json:"unit_id" validate:"required"`
	StudentId     string `json:"student_id" validate:"required"`
	StudentUnitId string `json:"student_unit_id" validate:"required"`
}

type StudentAttendanceLog struct {
	Date       string `json:"date"`
	LoginTime  string `json:"login_time"`
	LogoutTime string `json:"logout_time"`
}

type PdfDownloadRequest struct {
	UnitId    string `json:"unit_id" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type Student struct {
	StudentId     string `json:"student_id"`
	StudentUnitId string `json:"student_unit_id"`
	StudentName   string `json:"student_name"`
	StudentUsn    string `json:"student_usn"`
	Department    string `json:"department"`
}

type StudentInterface interface {
	CreateNewStudent(r *http.Request) error
	UpdateStudentDetails(r *http.Request) error
	DeleteStudent(r *http.Request) error
	GetStudentDetails(r *http.Request) ([]*Student, error)
	GetStudentLogs(r *http.Request) ([]*StudentAttendanceLog, error)
	DownloadPdf(r *http.Request) error
}
