package utils

import (
	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/xuri/excelize/v2"
)

func SetExcelHeader(file *excelize.File, sheetName string, dates []int) error {
	headers := make([]int, 0)

	headers = append(headers, dates...)

	cell, err := excelize.CoordinatesToCellName(1, 1)

	if err != nil {
		return err
	}

	if err := file.SetCellValue(sheetName, cell, "USN"); err != nil {
		return err
	}

	cell, err = excelize.CoordinatesToCellName(2, 1)

	if err != nil {
		return err
	}

	if err := file.SetCellValue(sheetName, cell, "Name"); err != nil {
		return err
	}

	for col, header := range headers {
		cell, err := excelize.CoordinatesToCellName(col+3, 1)

		if err != nil {
			return err
		}

		if err := file.SetCellValue(sheetName, cell, header); err != nil {
			return err
		}
	}

	return nil
}

func AddStudentsToExcelFile(file *excelize.File, sheetName string, students []*models.ExcelStudentInfo) error {
	for row, student := range students {
		cell, err := excelize.CoordinatesToCellName(1, row+2)

		if err != nil {
			return err
		}

		if err := file.SetCellValue(sheetName, cell, student.Usn); err != nil {
			return err
		}

		cell, err = excelize.CoordinatesToCellName(2, row+2)

		if err != nil {
			return err
		}

		if err := file.SetCellValue(sheetName, cell, student.Name); err != nil {
			return err
		}
	}

	return nil
}

func AddStudentAttendanceStatusToExcelFile(file *excelize.File, sheetName string, column int, studentAttendanceStatus []*models.ExcelStudentAttendanceStatus) error {
	for row, studentAttendanceStatus := range studentAttendanceStatus {
		cell, err := excelize.CoordinatesToCellName(column, row+2)

		if err != nil {
			return err
		}

		if err := file.SetCellValue(sheetName, cell, studentAttendanceStatus.AttendanceStatus); err != nil {
			return err
		}
	}
	return nil
}
