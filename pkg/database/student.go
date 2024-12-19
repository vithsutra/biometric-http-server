package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

func (q *Query) GenerateStudentAttendenceReport(unitId string, date string) ([]models.Student, error) {
	query := fmt.Sprintf(`SELECT 
    s.student_name,
    s.student_usn,
    COALESCE(a.login, 'Absent') AS login_time,
    COALESCE(a.logout, 'Absent') AS logout_time
FROM 
    %s s
LEFT JOIN 
    attendance a 
ON 
    s.student_id = a.student_id AND a.date = $1
ORDER BY 
    s.student_name;
`, unitId)
	res, err := q.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var student models.Student
	var students []models.Student
	for res.Next() {
		if err := res.Scan(&student.StudentName, &student.StudentUSN, &student.LoginTime, &student.LogoutTime); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	if res.Err() != nil {
		return nil, err
	}
	return students, nil
}

func (q *Query) NewStudent(student models.Student) error {
	tx , err := q.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}else{
			tx.Commit()
		}
	} ()
	if err != nil {
		return err
	}
	_ , err = tx.Exec("INSERT INTO fingerprintdata(student_id , student_unit_id , unit_id , fingerprint) VALUES($1 , $2 , $3 , $4)" , student.StudentId , student.StudentUnitId , student.UnitId , student.FingerprintData)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s(student_id , student_unit_id , student_name , student_usn , department) VALUES($1 , $2 , $3 , $4 , $5)" , student.UnitId)
	_ , err = tx.Exec(query , student.StudentId , student.StudentUnitId , student.StudentName , student.StudentUSN , student.Department)
	if err != nil {
		return err
	}
	_ , err = tx.Exec("INSERT INTO inserts(unit_id , student_unit_id , fingerprint_data) VALUES($1,$2,$3)" , student.UnitId , student.StudentUnitId , student.FingerprintData)
	if err != nil {
		return err
	}
	return nil
}

func (q *Query) FetchStudentDetails(unitId string) ([]models.Student,error) {
	query := fmt.Sprintf("SELECT student_id , student_name , student_usn , student_unit_id , department FROM %s" , unitId)
	res , err := q.db.Query(query)
	if err != nil {
		return nil,err
	}
	defer res.Close()
	var studentDetails models.Student
	var studentsDetails []models.Student
	for res.Next() {
		if err := res.Scan(&studentDetails.StudentId , &studentDetails.StudentUSN , &studentDetails.StudentUnitId , &studentDetails.Department); err != nil {
			return nil,err
		}
		studentsDetails = append(studentsDetails, studentDetails)
	}
	return studentsDetails , nil
}

func (q *Query) FetchStudentLogHistory(studentId string) ([]models.Student , error) {
	res , err := q.db.Query("SELECT login , logout , date FROM attendance WHERE student_id=$1" , studentId)
	if err != nil {
		return nil , err
	}
	defer res.Close()
	var log models.Student
	var logs []models.Student
	for res.Next() {
		if err := res.Scan(&log.LoginTime , &log.LogoutTime , &log.Date); err != nil {
			return nil , err
		}
		logs = append(logs, log)
	}
	if res.Err() != nil {
		return nil,err
	}
	return logs , nil
}

func (q *Query) DeleteStudent(unitid , studentUnitId , studentid string) error {
	tx , err := q.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}else{
			tx.Commit()
		}
	} ()
	if err != nil {
		return err
	}
	_ , err = tx.Exec("DELETE FROM fingerprintdata WHERE student_id=$1" , studentid)
	if err != nil {
		return err
	}
	_ , err = tx.Exec("INSERT INTO deletes(unit_id , student_unit_id) VALUES($1,$2)" , unitid , studentUnitId)
	if err != nil {
		return err
	}
	return nil
}

func (q *Query) UpdateStudent(student models.Student) error {
	query := fmt.Sprintf("UPDATE %s SET student_name=$1 , student_usn=$2 , department=$3 WHERE student_id=$4" , student.UnitId)
	_ , err := q.db.Exec(query , student.StudentName , student.StudentUSN , student.Department , student.StudentId)
	if err != nil {
		return err
	}
	return nil
}