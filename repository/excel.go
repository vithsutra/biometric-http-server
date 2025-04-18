package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/xuri/excelize/v2"
)

type ExcelRepository struct {
	DB *sql.DB
}

func NewExcelRepository(db *sql.DB) *ExcelRepository {
	return &ExcelRepository{DB: db}
}
func (r *ExcelRepository) DownloadExcel(req *models.ExcelDownloadRequest) (*excelize.File, error) {
	// Parse date range
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	// Step 1: Get distinct student info
	studentsQuery := `
	SELECT DISTINCT f.student_unit_id AS usn, COALESCE(s.student_name, 'N/A') AS name
	FROM fingerprintdata f
	LEFT JOIN ` + req.UnitId + ` s ON s.student_unit_id = f.student_unit_id
	WHERE f.unit_id = $1
	ORDER BY f.student_unit_id;
	`

	studentsRows, err := r.DB.Query(studentsQuery, req.UnitId)
	if err != nil {
		return nil, err
	}
	defer studentsRows.Close()

	type student struct {
		USN  string
		Name string
	}
	var students []student

	for studentsRows.Next() {
		var s student
		if err := studentsRows.Scan(&s.USN, &s.Name); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	// Step 2: Create Excel file
	file := excelize.NewFile()
	sheet := file.GetSheetName(0)

	// Step 3: Build headers (Name, USN, Date1, Date2, ...)
	headers := []string{"Name", "USN"}
	dateMap := make(map[int]string) // index -> date string

	// Adding dates in headers
	colIndex := 3
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		headers = append(headers, dateStr)
		dateMap[colIndex] = dateStr
		colIndex++
	}

	// Set header cells
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheet, cell, h)
	}

	// Step 4: Fill each student row
	row := 2
	for _, stu := range students {
		file.SetCellValue(sheet, fmt.Sprintf("A%d", row), stu.Name)
		file.SetCellValue(sheet, fmt.Sprintf("B%d", row), stu.USN)

		// Step 5: Query attendance per day and fill attendance columns
		for col, date := range dateMap {
			var status string
			query := `
			SELECT CASE
				WHEN login IS NOT NULL THEN 'Present'
				ELSE 'Absent'
			END AS status
			FROM attendance
			WHERE student_id = (SELECT student_id FROM fingerprintdata WHERE student_unit_id = $1 AND unit_id = $2 LIMIT 1)
			AND unit_id = $2 AND date = $3
			LIMIT 1;
			`
			err := r.DB.QueryRow(query, stu.USN, req.UnitId, date).Scan(&status)
			if err != nil {
				status = "Absent" // if no record found, mark as Absent
			}
			cell, _ := excelize.CoordinatesToCellName(col, row)
			file.SetCellValue(sheet, cell, status)
		}
		row++
	}

	return file, nil
}
