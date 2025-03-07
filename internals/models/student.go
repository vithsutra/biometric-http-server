package models

import (
	"net/http"

	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
)

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
	UserId    string `json:"user_id" validate:"required"`
	Slot      string `json:"slot" validate:"required,slot"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type ExcelDownloadRequest struct {
	UnitId    string `json:"unit_id" validate:"required"`
	UserId    string `json:"user_id" validate:"required"`
	Slot      string `json:"slot" validate:"required,slot"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type UserTime struct {
	MorningStart   string
	MorningEnd     string
	AfterNoonStart string
	AfterNoonEnd   string
	EveningStart   string
	EveningEnd     string
}

type StudentForPdf struct {
	StudentId string
	Name      string
	Usn       string
}

type PdfFormat struct {
	StudentId string `json:"student_id"`
	Usn       string `json:"usn"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Logout    string `json:"logout"`
}

type ExcelStudentInfo struct {
	Usn  string
	Name string
}

type ExcelStudentAttendanceStatus struct {
	Usn              string
	AttendanceStatus string
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
	DownloadPdf(r *http.Request) (*gopdf.GoPdf, error)
	DownloadExcel(r *http.Request) (*excelize.File, error)
}
