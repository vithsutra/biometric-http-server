package models

import (
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadRequest struct {
	UnitId    string `json:"unit_id" validate:"required"`
	UserId    string `json:"user_id" validate:"required"`
	Slot      string `json:"slot" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

type ExcelStudentInfo struct {
	Usn  string
	Name string
}

type ExcelStudentAttendanceStatus struct {
	Usn              string
	AttendanceStatus string
}

type ExcelInterface interface {
	DownloadExcel(req *ExcelDownloadRequest) (*excelize.File, error)
}
