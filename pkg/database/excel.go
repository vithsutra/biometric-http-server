package database

import "github.com/VsenseTechnologies/biometric_http_server/internals/models"

func(q *Query) GenerateExcelReport(unitid , start , end string) ([]models.Excel , error) {
	query := `
		SELECT 
    s.student_name, 
    s.student_usn, 
    TO_CHAR(a.date::date, 'DD') AS day, 
    CASE 
        WHEN a.login BETWEEN t.morning_start AND t.morning_end 
             AND a.logout BETWEEN t.evening_start AND t.evening_end THEN 'P'
        WHEN a.login BETWEEN t.morning_start AND t.morning_end 
             AND a.logout BETWEEN t.afternoon_start AND t.afternoon_end THEN 'MP'
        WHEN a.login BETWEEN t.afternoon_start AND t.afternoon_end 
             AND a.logout BETWEEN t.evening_start AND t.evening_end THEN 'AP'
        ELSE 'A'
    END AS status
FROM 
    attendance a
JOIN 
    fingerprintdata s ON a.student_id = s.student_id 
JOIN 
    biometric b ON a.unit_id = b.unit_id
JOIN 
    times t ON b.user_id = t.user_id 
WHERE 
    a.date BETWEEN $1 AND $2 
    AND a.unit_id = $3 
ORDER BY 
    s.student_usn, a.date;

	`;
	rows , err := q.db.Query(query , start , end , unitid )
    if err != nil {
        return nil , err
    }
    defer rows.Close()
    var attendance models.Excel
    var attendances []models.Excel
    for rows.Next() {
        if err := rows.Scan(&attendance.StudentName , &attendance.StudentUsn , &attendance.Date , &attendance.Status); err != nil {
            return nil , err
        }
        attendances = append(attendances, attendance)
    }
    if rows.Err() != nil {
        return nil , err
    }
    return attendances , nil
}