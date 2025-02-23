package database

// import (
// 	"fmt"
// 	"log"

// 	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
// )

// func(q *Query) GenerateExcelReport(unitid , start , end string) ([]models.Excel , error) {
// 	query := fmt.Sprintf( `
// 		WITH date_range AS (
//     SELECT generate_series($1::DATE, $2::DATE, INTERVAL '1 day')::DATE AS date
// ),
// students_with_dates AS (
//     SELECT
//         d.date,
//         u.student_name,
//         u.student_usn,
//         u.student_id,
//         b.unit_id  -- Get the unit_id from the biometric table
//     FROM
//         date_range d
//     CROSS JOIN
//         %s u
//     JOIN
//         biometric b ON u.student_id = b.user_id  -- Ensure the student is mapped to the correct unit_id
//     WHERE
//         b.unit_id = $3  -- Filter by the specific unit_id
// )
// SELECT
//     s.student_name,
//     s.student_usn,
//     TO_CHAR(s.date, 'DD') AS day,
//     COALESCE(
//         CASE
//             WHEN a.login BETWEEN t.morning_start AND t.morning_end
//                  AND a.logout BETWEEN t.evening_start AND t.evening_end THEN 'P'
//             WHEN a.login BETWEEN t.morning_start AND t.morning_end
//                  AND a.logout BETWEEN t.afternoon_start AND t.afternoon_end THEN 'MP'
//             WHEN a.login BETWEEN t.afternoon_start AND t.afternoon_end
//                  AND a.logout BETWEEN t.evening_start AND t.evening_end THEN 'AP'
//             ELSE 'A'
//         END, 'A'
//     ) AS status
// FROM
//     students_with_dates s
// LEFT JOIN
//     attendance a ON s.student_id = a.student_id AND s.date = CAST(a.date AS DATE)
// LEFT JOIN
//     biometric b ON a.unit_id = b.unit_id
// LEFT JOIN
//     times t ON b.user_id = t.user_id
// ORDER BY
//     s.student_usn, s.date;

// 	` , unitid);
// 	rows , err := q.db.Query(query , start , end , unitid )
//     if err != nil {
//         return nil , err
//     }
//     defer rows.Close()
//     var attendance models.Excel
//     var attendances []models.Excel
//     for rows.Next() {
//         if err := rows.Scan(&attendance.StudentName , &attendance.StudentUsn , &attendance.Date , &attendance.Status); err != nil {
//             return nil , err
//         }
//         attendances = append(attendances, attendance)
//     }
//     if rows.Err() != nil {
//         return nil , err
//     }
//     log.Println(attendances)
//     return attendances , nil
// }
