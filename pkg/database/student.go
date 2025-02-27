package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)

// 	_ , err = tx.Exec("INSERT INTO fingerprintdata(student_id , student_unit_id , unit_id , fingerprint) VALUES($1 , $2 , $3 , $4)" , student.StudentId , student.StudentUnitId , student.UnitId , student.FingerprintData)
// 	if err != nil {
// 		return err
// 	}
// 	query := fmt.Sprintf("INSERT INTO %s(student_id , student_unit_id , student_name , student_usn , department) VALUES($1 , $2 , $3 , $4 , $5)" , student.UnitId)
// 	_ , err = tx.Exec(query , student.StudentId , student.StudentUnitId , student.StudentName , student.StudentUSN , student.Department)
// 	if err != nil {
// 		return err
// 	}
// 	_ , err = tx.Exec("INSERT INTO inserts(unit_id , student_unit_id , fingerprint_data) VALUES($1,$2,$3)" , student.UnitId , student.StudentUnitId , student.FingerprintData)
// 	if err != nil {
// 		return err

func (q *Query) CreateNewStudent(student *models.Student, unitId string, fingerPrintData string) error {
	query1 := `INSERT INTO fingerprintdata (student_id,student_unit_id,unit_id,fingerprint) VALUES ($1,$2,$3,$4)`
	query2 := fmt.Sprintf("INSERT INTO %s (student_id,student_unit_id,student_name,student_usn,department) VALUES ($1,$2,$3,$4,$5)", unitId)
	query3 := `INSERT INTO inserts (unit_id,student_unit_id,fingerprint_data) VALUES ($1,$2,$3)`

	tx, err := q.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.Exec(query1, student.StudentId, student.StudentUnitId, unitId, fingerPrintData); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query2, student.StudentId, student.StudentUnitId, student.StudentName, student.StudentUsn, student.Department); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query3, unitId, student.StudentUnitId, fingerPrintData); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (q *Query) UpdateStudent(unitId string, studentId string, studentName string, studentUsn string, department string) error {
	query := fmt.Sprintf(`UPDATE %s SET student_name=$2,student_usn=$3,department=$4 WHERE student_id=$1`, unitId)
	if _, err := q.db.Exec(query, studentId, studentName, studentUsn, department); err != nil {
		return err
	}
	return nil
}

func (q *Query) DeleteStudent(unitId string, studentId string, studentUnitId string) error {
	query1 := `DELETE FROM fingerprintdata WHERE student_id=$1`
	query2 := `INSERT INTO deletes (unit_id,student_unit_id) VALUES ($1,$2)`
	query3 := `DELETE FROM inserts WHERE unit_id=$1 AND student_unit_id=$2`

	tx, err := q.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.Exec(query1, studentId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query2, unitId, studentUnitId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query3, unitId, studentUnitId); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (q *Query) GetStudentDetails(unitId string) ([]*models.Student, error) {
	query := fmt.Sprintf(`SELECT student_id,student_unit_id,student_name,student_usn,department FROM %s`, unitId)

	var students []*models.Student

	rows, err := q.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var student models.Student

		if err := rows.Scan(&student.StudentId, &student.StudentUnitId, &student.StudentName, &student.StudentUsn, &student.Department); err != nil {
			return nil, err
		}

		students = append(students, &student)
	}

	return students, nil
}

func (q *Query) GetStudentLogs(studentId string) ([]*models.StudentAttendanceLog, error) {
	query := `SELECT date,login,logout FROM attendance WHERE student_id=$1`

	rows, err := q.db.Query(query, studentId)

	var studentLogs []*models.StudentAttendanceLog

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var attendanceLog models.StudentAttendanceLog

		if err := rows.Scan(&attendanceLog.Date, &attendanceLog.LoginTime, &attendanceLog.LogoutTime); err != nil {
			return nil, err
		}

		if attendanceLog.LoginTime != "pending" {
			t1, err := utils.ConvertTo12HourFormat(attendanceLog.LoginTime)

			if err != nil {
				return nil, err
			}

			attendanceLog.LoginTime = t1

		}

		if attendanceLog.LogoutTime != "pending" {
			t2, err := utils.ConvertTo12HourFormat(attendanceLog.LogoutTime)

			if err != nil {
				return nil, err
			}
			attendanceLog.LogoutTime = t2
		}
		studentLogs = append(studentLogs, &attendanceLog)
	}

	return studentLogs, nil
}

func (q *Query) GetStudentsAttendanceLogForPdf(unitId string, userId string, date string, slot string) ([]*models.PdfFormat, error) {
	query := `	
	WITH student_logs AS (
		SELECT 
			s.student_name,
			s.student_usn,
			a.student_id,
			a.date,
			COALESCE(a.login, 'Pending') AS login,
			COALESCE(a.logout, 'Pending') AS logout
		FROM ` + unitId + ` s
		LEFT JOIN attendance a 
			ON s.student_id = a.student_id 
			AND s.student_unit_id = a.student_unit_id
	)
	SELECT 
		sl.student_name,
		sl.student_usn,
		COALESCE(a.login, 'Pending') AS login,
		COALESCE(a.logout, 'Pending') AS logout
	FROM student_logs sl
	LEFT JOIN times t 
		ON t.user_id = $1  -- Filter times for the given user_id
	LEFT JOIN attendance a 
		ON sl.student_id = a.student_id 
		AND sl.date = $2
	WHERE 
		sl.date = $2
		AND (
			($3 = 'morning'   AND a.login::time >= t.morning_start::time  AND a.logout::time <= t.afternoon_end::time)
			OR ($3 = 'afternoon' AND a.login::time >= t.afternoon_start::time AND a.logout::time <= t.evening_end::time)
			OR ($3 = 'full'    AND a.login::time >= t.morning_start::time  AND a.logout::time <= t.evening_end::time)
		);`

	rows, err := q.db.Query(query, userId, date, slot)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var studentLogs []*models.PdfFormat

	for rows.Next() {
		var studentLog models.PdfFormat

		if err := rows.Scan(&studentLog.Name, &studentLog.Usn, &studentLog.Login, &studentLog.Logout); err != nil {
			return nil, err
		}

		studentLogs = append(studentLogs, &studentLog)
	}

	return studentLogs, nil

}
