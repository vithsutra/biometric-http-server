package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

func (q *Query) GenerateStudentAttendenceReport(unitId string) ([]models.Student,error){
	query := fmt.Sprintf("SELECT student_name , student_usn , login , logout FROM %s",unitId)
	res , err := q.db.Query(query)
	if err != nil {
		return nil,err
	}
	defer res.Close()
	var student models.Student
	var students []models.Student
	for res.Next() {
		if err := res.Scan(&student.StudentName , &student.StudentUSN , &student.LoginTime , &student.LogoutTime); err != nil {
			return nil,err
		}
		students = append(students, student)
	}
	if res.Err() != nil {
		return nil , err
	}
	return students , nil
}