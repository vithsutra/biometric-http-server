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

// Helper to parse either "02-01-2006" or RFC3339
func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339, // "2025-03-22T09:00:00Z"
		"02-01-2006", // "22-03-2025"
		"2006-01-02", // "2025-03-22"
		"02/01/2006", // "22/03/2025"
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}

func (r *ExcelRepository) DownloadExcel(req *models.ExcelDownloadRequest) (*excelize.File, error) {
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %v", err)
	}

	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %v", err)
	}

	studentsQuery := `
	SELECT DISTINCT f.student_unit_id AS usn, COALESCE(s.student_name, 'N/A') AS name,
	s.student_name
	FROM fingerprintdata f
	LEFT JOIN ` + req.UnitId + ` s ON s.student_unit_id = f.student_unit_id
	WHERE f.unit_id = $1
	ORDER BY s.student_name;
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
		var rawName string
		if err := studentsRows.Scan(&s.USN, &s.Name, &rawName); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	file := excelize.NewFile()
	sheet := file.GetSheetName(0)

	headers := []string{"Name", "USN"}
	dateMap := make(map[int]string)

	colIndex := 3
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dbDateStr := d.Format("2006-01-02")
		displayDateStr := d.Format("02/01/2006")

		headers = append(headers, displayDateStr)
		dateMap[colIndex] = dbDateStr
		colIndex++
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheet, cell, h)
	}

	row := 2
	for _, stu := range students {
		file.SetCellValue(sheet, fmt.Sprintf("A%d", row), stu.Name)
		file.SetCellValue(sheet, fmt.Sprintf("B%d", row), stu.USN)

		for col, date := range dateMap {
			var status string
			query := `
			SELECT CASE
				WHEN login IS NOT NULL THEN 'P'
				ELSE 'A'
			END AS status
			FROM attendance
			WHERE student_id = (
				SELECT student_id FROM fingerprintdata 
				WHERE student_unit_id = $1 AND unit_id = $2 LIMIT 1
			)
			AND unit_id = $2 AND date = $3
			LIMIT 1;
			`
			err := r.DB.QueryRow(query, stu.USN, req.UnitId, date).Scan(&status)
			if err != nil {
				status = "A"
			}
			cell, _ := excelize.CoordinatesToCellName(col, row)
			file.SetCellValue(sheet, cell, status)
		}
		row++
	}

	headerStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#000000"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	centerStyle, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})

	file.SetCellStyle(sheet, "A1", fmt.Sprintf("%s1", string(rune('A'+len(headers)-1))), headerStyle)

	lastCol, _ := excelize.ColumnNumberToName(len(headers))
	file.SetCellStyle(sheet, "A2", fmt.Sprintf("%s%d", lastCol, row-1), centerStyle)

	file.SetColWidth(sheet, "A", "A", 25)
	file.SetColWidth(sheet, "B", "B", 20)
	file.SetColWidth(sheet, "C", lastCol, 18)

	return file, nil
}
