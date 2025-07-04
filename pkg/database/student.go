package database

import (

	// "regexp"
	"sync"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/lib/pq"
)

//	func isValidIdentifier(id string) bool {
//		re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
//		return re.MatchString(id)
//	}
func (q *Query) CheckStudentUnitIdExists(unitId string, studentUnitId string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM  ` + unitId + ` WHERE student_unit_id = $1)`

	var isStudentUnitIdExists bool
	err := q.db.QueryRow(query, studentUnitId).Scan(&isStudentUnitIdExists)
	return isStudentUnitIdExists, err
}

func (q *Query) CreateNewStudent(student *models.Student, unitId string, fingerPrintData []string) error {
	query1 := "INSERT INTO student (student_id,unit_id,student_name,student_usn,department) VALUES ($1,$2,$3,$4,$5)"
	query2 := `INSERT INTO fingerprintdata (student_id,student_unit_id,unit_id,fingerprint) VALUES ($1,$2,$3,$4)`
	query3 := `INSERT INTO inserts (unit_id,student_unit_id,fingerprint_data) VALUES ($1,$2,$3)`

	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(query1, student.StudentId, unitId, student.StudentName, student.StudentUsn, student.Department); err != nil {
		tx.Rollback()
		return err
	}

	for i := 0; i < 6; i++ {
		if _, err := tx.Exec(query2, student.StudentId, student.StudentUnitId[i], unitId, fingerPrintData[i]); err != nil {
			tx.Rollback()
			return err
		}
	}

	for i := 0; i < 6; i++ {
		if _, err := tx.Exec(query3, unitId, student.StudentUnitId[i], fingerPrintData[i]); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := q.UpdateAvailableStudentUnitIDs(unitId, student.StudentUnitId, false); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (q *Query) UpdateStudent(unitId string, studentId string, studentName string, studentUsn string, department string) error {
	query := `UPDATE student SET student_name=$2,student_usn=$3,department=$4 WHERE student_id=$1`
	_, err := q.db.Exec(query, studentId, studentName, studentUsn, department)
	return err
}

func (q *Query) DeleteStudent(unitId string, studentId string) error {
	query1 := `DELETE FROM fingerprintdata WHERE student_id=$1 RETURNING student_unit_id;`
	query2 := `INSERT INTO deletes (unit_id,student_unit_id) VALUES ($1,$2)`
	query3 := `DELETE FROM inserts WHERE unit_id=$1 AND student_unit_id=$2`
	query4 := `DELETE FROM student WHERE student_id=$1 AND unit_id=$2`

	tx, err := q.db.Begin()

	if err != nil {
		return err
	}

	rows, err := tx.Query(query1, studentId)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	var student_unit_ids []string
	var student_unit_id string

	for rows.Next() {
		if err := rows.Scan(&student_unit_id); err != nil {
			tx.Rollback()
			return err
		}
		student_unit_ids = append(student_unit_ids, student_unit_id)
	}

	for i := 0; i < 6; i++ {
		if _, err := tx.Exec(query2, unitId, student_unit_ids[i]); err != nil {
			tx.Rollback()
			return err
		}
	}
	for i := 0; i < 6; i++ {
		if _, err := tx.Exec(query3, unitId, student_unit_ids[i]); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := q.UpdateAvailableStudentUnitIDs(unitId, student_unit_ids, true); err != nil {
		tx.Rollback()
		return err
	}

	if rows.Err() != nil {
		tx.Rollback()
		return rows.Err()
	}

	if _, err := tx.Exec(query4, studentId, unitId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// func (q *Query) GetStudentDetails(unitId string, limit, offset int) ([]*models.Student, int, error) {
// 	query1 := `SELECT
// 				s.student_id,
// 				s.student_name,
// 				s.student_usn,
// 				s.department,
// 				ARRAY_AGG(fd.student_unit_id ORDER BY fd.student_unit_id) AS student_unit_id
// 			FROM
// 				student s
// 			JOIN
// 				fingerprintdata fd ON s.student_id = fd.student_id
// 			WHERE
// 				s.unit_id = $1
// 			GROUP BY
// 				s.student_id, s.student_name, s.student_usn, s.department
// 			HAVING
// 				COUNT(fd.student_unit_id) = 6
// 			LIMIT $2
// 			OFFSET $3;
// 			`

// 	query2 := `SELECT
// 					COUNT(*) AS total_students
// 				FROM (
// 					SELECT s.student_id
// 					FROM student s
// 					WHERE s.unit_id = $1
// 					GROUP BY s.student_id
// 				) AS sub;`
// 	var students []*models.Student

// 	tx, err := q.db.Begin()
// 	if err != nil {
// 		return nil, -1, err
// 	}
// 	rows, err := tx.Query(query1, unitId, limit, offset)
// 	if err != nil {
// 		tx.Rollback()
// 		return nil, -1, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var student models.Student

// 		if err := rows.Scan(&student.StudentId, &student.StudentName, &student.StudentUsn, &student.Department, pq.Array(&student.StudentUnitId)); err != nil {
// 			tx.Rollback()
// 			return nil, -1, err
// 		}

// 		students = append(students, &student)
// 	}

// 	if rows.Err() != nil {
// 		tx.Rollback()
// 		return nil, -1, rows.Err()
// 	}

// 	var totalStudents int

// 	if err := tx.QueryRow(query2, unitId).Scan(&totalStudents); err != nil {
// 		tx.Rollback()
// 		return nil, -1, err
// 	}

// 	tx.Commit()
// 	return students, totalStudents, nil
// }

func (q *Query) GetStudentDetails(unitId string) ([]*models.Student, error) {
	query1 := `SELECT 
		s.student_id,
		s.student_name,
		s.student_usn,
		s.department,
		ARRAY_AGG(fd.student_unit_id ORDER BY fd.student_unit_id) AS student_unit_id
		FROM 
			student s
		JOIN 
			fingerprintdata fd ON s.student_id = fd.student_id
		WHERE
			s.unit_id = $1
		GROUP BY 
			s.student_id, s.student_name, s.student_usn, s.department
		HAVING 
			COUNT(fd.student_unit_id) = 6
	`
	rows, err := q.db.Query(query1, unitId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []*models.Student

	for rows.Next() {
		var student models.Student

		if err := rows.Scan(&student.StudentId, &student.StudentName, &student.StudentUsn, &student.Department, pq.Array(&student.StudentUnitId)); err != nil {
			return nil, err
		}

		students = append(students, &student)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
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

	defer rows.Close()

	for rows.Next() {
		var attendanceLog models.StudentAttendanceLog

		if err := rows.Scan(&attendanceLog.Date, &attendanceLog.LoginTime, &attendanceLog.LogoutTime); err != nil {
			return nil, err
		}

		if attendanceLog.LoginTime != "25:00" {
			t1, err := utils.ConvertTo12HourFormat(attendanceLog.LoginTime)

			if err != nil {
				return nil, err
			}

			attendanceLog.LoginTime = t1

		}

		if attendanceLog.LogoutTime != "25:00" {
			t2, err := utils.ConvertTo12HourFormat(attendanceLog.LogoutTime)

			if err != nil {
				return nil, err
			}
			attendanceLog.LogoutTime = t2
		}
		studentLogs = append(studentLogs, &attendanceLog)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return studentLogs, nil
}

func (q *Query) GetStudentsCountFromUnit(unitId string) (int32, error) {
	query := `SELECT COUNT(*) FROM student WHERE unit_id = $1`

	var studentCount int32
	err := q.db.QueryRow(query, unitId).Scan(&studentCount)
	return studentCount, err

}

func (q *Query) GetUserStandardTime(userId string) (*models.UserTime, error) {
	query := `SELECT morning_start,morning_end,afternoon_start,afternoon_end,evening_start,evening_end FROM times where user_id=$1`
	var userTime models.UserTime
	if err := q.db.QueryRow(query, userId).Scan(
		&userTime.MorningStart,
		&userTime.MorningEnd,
		&userTime.AfterNoonStart,
		&userTime.AfterNoonEnd,
		&userTime.EveningStart,
		&userTime.EveningEnd,
	); err != nil {
		return nil, err
	}
	return &userTime, nil
}

func (q *Query) GetStudentsForPdf(unitId string, studentsCount int32) (map[string]*models.PdfFormat, error) {
	query := `SELECT student_id, student_name, student_usn FROM student WHERE unit_id = $1 ORDER BY student_usn`

	rows, err := q.db.Query(query, unitId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pdfFormats = make(map[string]*models.PdfFormat, studentsCount)

	for rows.Next() {
		var pdfFormat models.PdfFormat
		if err := rows.Scan(&pdfFormat.StudentId, &pdfFormat.Name, &pdfFormat.Usn); err != nil {
			return nil, err
		}

		pdfFormats[pdfFormat.StudentId] = &pdfFormat
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return pdfFormats, nil
}

func (q *Query) GetStudentsAttendanceLogForPdf(studentsCount int32, userTime *models.UserTime, pdfFormats map[string]*models.PdfFormat, date string, slot string) error {

	var wg sync.WaitGroup

	for studentId := range pdfFormats {

		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			query := `SELECT login,logout FROM attendance WHERE date=$1 AND student_id=$2`
			rows, err := q.db.Query(query, date, studentId)
			if err != nil {
				return
			}

			defer rows.Close()

			var isStudentEntryValid bool

			for rows.Next() {

				var studentAttendanceLog models.StudentAttendanceLog

				if err := rows.Scan(&studentAttendanceLog.LoginTime, &studentAttendanceLog.LogoutTime); err != nil {
					return
				}

				if studentAttendanceLog.LogoutTime == "25:00" {
					continue
				}

				if slot == "morning" {

					entryValid, err := utils.CompareWithStandardTime(userTime.MorningStart, userTime.AfterNoonEnd, studentAttendanceLog.LoginTime, studentAttendanceLog.LogoutTime)

					if err != nil {
						return
					}

					if entryValid {
						isStudentEntryValid = true
						pdfFormats[studentId].Login = studentAttendanceLog.LoginTime
						pdfFormats[studentId].Logout = studentAttendanceLog.LogoutTime
						break
					}

				} else if slot == "evening" {
					entryValid, err := utils.CompareWithStandardTime(userTime.AfterNoonStart, userTime.EveningEnd, studentAttendanceLog.LoginTime, studentAttendanceLog.LogoutTime)

					if err != nil {
						return
					}

					if entryValid {
						isStudentEntryValid = true
						pdfFormats[studentId].Login = studentAttendanceLog.LoginTime
						pdfFormats[studentId].Logout = studentAttendanceLog.LogoutTime
						break
					}

				} else {
					entryValid, err := utils.CompareWithStandardTime(userTime.MorningStart, userTime.EveningEnd, studentAttendanceLog.LoginTime, studentAttendanceLog.LogoutTime)

					if err != nil {
						return
					}

					if entryValid {
						isStudentEntryValid = true
						pdfFormats[studentId].Login = studentAttendanceLog.LoginTime
						pdfFormats[studentId].Logout = studentAttendanceLog.LogoutTime
						break
					}
				}
			}

			if !isStudentEntryValid {
				pdfFormats[studentId].Login = "pending"
				pdfFormats[studentId].Logout = "pending"
			}
		}(&wg)
	}

	wg.Wait()
	return nil
}
