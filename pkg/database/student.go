package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)

func (q *Query) CheckStudentUnitIdExists(unitId string, studentUnitId string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM  ` + unitId + ` WHERE student_unit_id = $1)`

	var isStudentUnitIdExists bool

	if err := q.db.QueryRow(query, studentUnitId).Scan(&isStudentUnitIdExists); err != nil {
		return false, err
	}

	return isStudentUnitIdExists, nil
}

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
	query := `
			SELECT date, login, logout 
				FROM attendance 
				WHERE student_id = $1 
			ORDER BY date DESC;
			`

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

	if slot == "morning" {
		query := `
		WITH student_list AS (
    SELECT s.student_id, s.student_name, s.student_usn
    FROM ` + unitId + ` s
), 
student_attendance AS (
    SELECT 
        a.student_id, 
        a.login::time AS login, 
        a.logout::time AS logout
    FROM attendance a
    WHERE a.date = $1  -- Replace with your provided date
),
time_reference AS (
    SELECT 
        morning_start::time AS morning_start, 
        afternoon_end::time AS afternoon_end 
    FROM times
    WHERE user_id = $2  -- Replace with your provided user_id
)
SELECT 
    sl.student_name,
    sl.student_usn,
    COALESCE(
        CASE 
            WHEN sa.logout = '00:00'::time THEN 'pending'
            WHEN sa.login >= tr.morning_start AND sa.logout <= tr.afternoon_end 
            THEN TO_CHAR(sa.login, 'HH24:MI')  -- Convert time to "hh:mm" format
            ELSE 'pending' 
        END, 
        'pending'
    ) AS login,
    COALESCE(
        CASE 
            WHEN sa.logout = '00:00'::time THEN 'pending'
            WHEN sa.login >= tr.morning_start AND sa.logout <= tr.afternoon_end 
            THEN TO_CHAR(sa.logout, 'HH24:MI')  -- Convert time to "hh:mm" format
            ELSE 'pending' 
        END, 
        'pending'
    ) AS logout
FROM student_list sl
LEFT JOIN student_attendance sa ON sl.student_id = sa.student_id
CROSS JOIN time_reference tr
ORDER BY sl.student_usn;
			`

		rows, err := q.db.Query(query, date, userId)

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

	if slot == "afternoon" {
		query := `
	WITH student_list AS (
    SELECT s.student_id, s.student_name, s.student_usn
    FROM ` + unitId + ` s
), 
student_attendance AS (
    SELECT 
        a.student_id, 
        a.login::time AS login, 
        a.logout::time AS logout
    FROM attendance a
    WHERE a.date = $1  -- Replace with your provided date
),
time_reference AS (
    SELECT 
        afternoon_start::time AS afternoon_start, 
        evening_end::time AS evening_end 
    FROM times
    WHERE user_id = $2  -- Replace with your provided user_id
)
SELECT 
    sl.student_name,
    sl.student_usn,
    COALESCE(
        CASE 
            WHEN sa.logout = '00:00'::time THEN 'pending'  
            WHEN sa.login >= tr.afternoon_start AND sa.logout <= tr.evening_end 
            THEN TO_CHAR(sa.login, 'HH24:MI')  -- Convert time to "hh:mm" format
            ELSE 'pending' 
        END, 
        'pending'
    ) AS login,
    COALESCE(
        CASE 
            WHEN sa.logout = '00:00'::time THEN 'pending'  
            WHEN sa.login >= tr.afternoon_start AND sa.logout <= tr.evening_end 
            THEN TO_CHAR(sa.logout, 'HH24:MI')  -- Convert time to "hh:mm" format
            ELSE 'pending' 
        END, 
        'pending'
    ) AS logout
FROM student_list sl
LEFT JOIN student_attendance sa ON sl.student_id = sa.student_id
CROSS JOIN time_reference tr
ORDER BY sl.student_usn;

	`

		rows, err := q.db.Query(query, date, userId)

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

	query := `
	WITH student_list AS (
    SELECT s.student_id, s.student_name, s.student_usn
    FROM ` + unitId + ` s
), 
student_attendance AS (
    SELECT 
        a.student_id, 
        a.login::time AS login, 
        a.logout::time AS logout
    FROM attendance a
    WHERE a.date = $1  -- Replace with your provided date
),
time_reference AS (
    SELECT 
        morning_start::time AS morning_start, 
        afternoon_end::time AS afternoon_end,
        afternoon_start::time AS afternoon_start, 
        evening_end::time AS evening_end 
    FROM times
    WHERE user_id = $2  -- Replace with your provided user_id
),
student_entries AS (
    SELECT 
        sa.student_id,
        -- Find the first login entry within the morning-afternoon range
        MAX(CASE 
            WHEN sa.login >= tr.morning_start AND sa.logout <= tr.afternoon_end 
            THEN sa.login 
        END) AS morning_login,
        -- Find the second logout entry within the afternoon-evening range
        MAX(CASE 
            WHEN sa.login >= tr.afternoon_start AND sa.logout <= tr.evening_end 
            THEN sa.logout 
        END) AS evening_logout
    FROM student_attendance sa
    CROSS JOIN time_reference tr
    GROUP BY sa.student_id
)
SELECT 
    sl.student_name,
    sl.student_usn,
    CASE 
        WHEN se.morning_login IS NOT NULL AND se.evening_logout IS NOT NULL 
        THEN TO_CHAR(se.morning_login, 'HH24:MI')  -- Convert time to "hh:mm" format
        ELSE 'pending' 
    END AS login,
    CASE 
        WHEN se.morning_login IS NOT NULL AND se.evening_logout IS NOT NULL 
        THEN TO_CHAR(se.evening_logout, 'HH24:MI')  -- Convert time to "hh:mm" format
        ELSE 'pending' 
    END AS logout
FROM student_list sl
LEFT JOIN student_entries se ON sl.student_id = se.student_id
ORDER BY sl.student_usn;

	`

	rows, err := q.db.Query(query, date, userId)

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
