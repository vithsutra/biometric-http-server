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
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

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
		dateStr := d.Format("02/01/2006")
		headers = append(headers, dateStr)
		dateMap[colIndex] = dateStr
		colIndex++
	}

	headerStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#000000"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheet, cell, h)
		file.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		width := float64(len(header) + 5)
		file.SetColWidth(sheet, colName, colName, width)
	}

	maxNameLength := len("Name")
	maxUSNLength := len("USN")

	for _, stu := range students {
		if len(stu.Name) > maxNameLength {
			maxNameLength = len(stu.Name)
		}
		if len(stu.USN) > maxUSNLength {
			maxUSNLength = len(stu.USN)
		}
	}

	file.SetColWidth(sheet, "A", "A", float64(maxNameLength+5))
	file.SetColWidth(sheet, "B", "B", float64(maxUSNLength+5))

	for col, date := range dateMap {
		maxDateLength := len(date)
		colName, _ := excelize.ColumnNumberToName(col)

		file.SetColWidth(sheet, colName, colName, float64(maxDateLength+5))
	}

	contentStyle, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	row := 2
	for _, stu := range students {

		file.SetCellValue(sheet, fmt.Sprintf("A%d", row), stu.Name)
		file.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), contentStyle)

		file.SetCellValue(sheet, fmt.Sprintf("B%d", row), stu.USN)
		file.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), contentStyle)

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
			file.SetCellStyle(sheet, cell, cell, contentStyle)
		}
		row++
	}

	return file, nil
}
