package database

func(q *Query) GenerateExcelReport() {
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
    fingerprintdata s ON a.student_id = s.student_id  -- Join fingerprintdata for student details
JOIN 
    biometric b ON a.unit_id = b.unit_id  -- Join biometric to get user_id
JOIN 
    times t ON b.user_id = t.user_id  -- Now join times using user_id from biometric
WHERE 
    a.date BETWEEN $1 AND $2 
    AND a.unit_id = $3  -- Condition to filter by unit_id
ORDER BY 
    s.student_usn, a.date;

	`;
	q.db.Query("")
}