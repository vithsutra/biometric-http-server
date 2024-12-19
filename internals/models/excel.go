package models

import (
	"net/http"

	"github.com/xuri/excelize/v2"
)

type Excel struct {
	StudentName string `json:"student_name"`
	StudentUsn  string `json:"student_usn"`
	Date        string `json:"day"`
	Status      string `json:"status"`
}

type ExcelDetails struct {
	UnitId    string `json:"unit_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ExcelInterface interface{
	GenerateExcelReport(*http.Request) (*excelize.File , string , error)
}
