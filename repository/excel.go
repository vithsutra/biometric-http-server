package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/xuri/excelize/v2"
)

type ExcelRepo struct {
	db *sql.DB
}

func NewExcelRepo(db *sql.DB) *ExcelRepo {
	return &ExcelRepo{
		db,
	}
}

func (er *ExcelRepo) GenerateExcelReport(r *http.Request) (*excelize.File , string , error) {
	var newExcel models.ExcelDetails
	if err := utils.Decode(r , &newExcel); err != nil {
		return nil , "" , err
	}
	query := database.NewQuery(er.db)
	attendance , err := query.GenerateExcelReport(newExcel.UnitId , newExcel.StartDate , newExcel.EndDate)
	if err != nil {
		return nil , "",err
	}
	startDate, err := time.Parse("2006-01-02", newExcel.StartDate)
	if err != nil {
		return nil , "" , err
	}
	endDate, err := time.Parse("2006-01-02", newExcel.EndDate)
	if err != nil {
		return nil , "" , err
	}

	// Generate date range
	dates := []string{}
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("02")) // Day of the month
	}

	f := excelize.NewFile()
	sheetName := "Attendance"
	f.NewSheet(sheetName)

	f.DeleteSheet("Sheet1")

	headers := append([]string{"Name", "USN"}, dates...)
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil , "" , err
		}
		style, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"000000"}, Pattern: 1},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		})
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "B", "B", 15)
	for i := 3; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		f.SetColWidth(sheetName, col, col, 5)
	}

	studentMap := make(map[string]map[string]string)
	for _, record := range attendance {
		if _, exists := studentMap[record.StudentUsn]; !exists {
			studentMap[record.StudentUsn] = map[string]string{
				"Name": record.StudentName,
			}
		}
		studentMap[record.StudentUsn][record.Date] = record.Status
	}

	row := 2
	for usn, data := range studentMap {
		if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), data["Name"]); err != nil {
			return nil , "" , err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), usn); err != nil {
			return nil , "" , err
		}

		for colIdx, day := range dates {
			status := data[day]                                      
			cell, _ := excelize.CoordinatesToCellName(colIdx+3, row)
			if err := f.SetCellValue(sheetName, cell, status); err != nil {
				return nil , "" , err
			}
		}
		row++
	}

	return f , newExcel.StartDate+"-"+newExcel.EndDate+"-Attendance.xlsx" , nil
}
