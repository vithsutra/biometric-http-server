package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
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

func (er *ExcelRepo) GenerateExcelReport(r *http.Request) (*excelize.File, string, error) {
	var newExcel models.ExcelDetails

	// Read and decode the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading request body: %w", err)
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return nil, "", fmt.Errorf("empty request body")
	}

	if err := json.Unmarshal(body, &newExcel); err != nil {
		return nil, "", fmt.Errorf("error decoding request body: %w", err)
	}

	// Fetch attendance records
	query := database.NewQuery(er.db)
	attendance, err := query.GenerateExcelReport(newExcel.UnitId, newExcel.StartDate, newExcel.EndDate)
	if err != nil {
		return nil, "", err
	}

	// Parse start and end dates
	startDate, err := time.Parse("2006-01-02", newExcel.StartDate)
	if err != nil {
		return nil, "", fmt.Errorf("invalid start date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", newExcel.EndDate)
	if err != nil {
		return nil, "", fmt.Errorf("invalid end date: %w", err)
	}

	// Generate date range for column headers
	dates := []string{}
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("02")) // Day of the month
	}

	// Create Excel file and sheet
	f := excelize.NewFile()
	sheetName := "Attendance"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	// Set headers
	headers := append([]string{"Name", "USN"}, dates...)
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, "", err
		}
		style, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"000000"}, Pattern: 1},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		})
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Adjust column widths
	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "B", "B", 15)
	for i := 3; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		f.SetColWidth(sheetName, col, col, 5)
	}

	// Map attendance records
	studentMap := make(map[string]map[string]string)
	for _, record := range attendance {
		if _, exists := studentMap[record.StudentUsn]; !exists {
			studentMap[record.StudentUsn] = map[string]string{
				"Name": record.StudentName,
			}
		}
		studentMap[record.StudentUsn][record.Date] = record.Status
	}

	// Populate rows with student data
	row := 2
	for usn, data := range studentMap {
		if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), data["Name"]); err != nil {
			return nil, "", err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), usn); err != nil {
			return nil, "", err
		}

		for colIdx, day := range dates {
			status := data[day]
			cell, _ := excelize.CoordinatesToCellName(colIdx+3, row)
			if err := f.SetCellValue(sheetName, cell, status); err != nil {
				return nil, "", err
			}
		}
		row++
	}

	// Return the Excel file
	return f, newExcel.StartDate + "-" + newExcel.EndDate + "-Attendance.xlsx", nil
}
