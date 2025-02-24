package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

func (q *Query) CreateUser(user *models.User) error {

	query1 := `INSERT INTO users (user_id,user_name,email,password) VALUES ($1,$2,$3,$4)`
	query2 := `INSERT INTO times (user_id,morning_start,morning_end,afternoon_start,afternoon_end,evening_start,evening_end) VALUES ($1,$2,$3,$4,$5,$6,$7)`

	tx, err := q.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.Exec(query1, user.UserId, user.UserName, user.Email, user.Password); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query2, user.UserId, `0001-01-01 00:00:00 +0000 UTC`, `0001-01-01 00:00:00 +0000 UTC`, `0001-01-01 00:00:00 +0000 UTC`, `0001-01-01 00:00:00 +0000 UTC`, `0001-01-01 00:00:00 +0000 UTC`, `0001-01-01 00:00:00 +0000 UTC`); err != nil {
		tx.Rollback()

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (q *Query) UserLogin(userName string) (string, string, error) {
	query := `SELECT user_id,password FROM users WHERE user_name = $1`

	var password string
	var userId string

	if err := q.db.QueryRow(query, userName).Scan(&userId, &password); err != nil {
		return "", "", err
	}

	return userId, password, nil
}

func (q *Query) UpdateTime(
	userId string,
	morningStartTime string,
	morningEndTime string,
	afterNoonStartTime string,
	afterNoonEndTime string,
	eveningStartTime string,
	eveningEndTime string,
) error {
	query := `UPDATE times SET morning_start=$2,morning_end=$3,afternoon_start=$4,afternoon_end=$5,evening_start=$6,evening_end=$7 WHERE user_id=$1`
	_, err := q.db.Exec(
		query,
		userId,
		morningStartTime,
		morningEndTime,
		afterNoonStartTime,
		afterNoonEndTime,
		eveningStartTime,
		eveningEndTime,
	)
	return err
}

func (q *Query) GiveUserAccess(userId string) (string, string, string, error) {
	query := `SELECT user_name,password,email FROM users WHERE user_id=$1`

	var userName string
	var password string
	var email string

	if err := q.db.QueryRow(query, userId).Scan(&userName, &password, &email); err != nil {
		return "", "", "", err
	}

	return userName, password, email, nil

}

func (q *Query) GetAllUsers() ([]*models.User, error) {
	query := `SELECT user_id,user_name,email FROM users`

	var users []*models.User

	rows, err := q.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.UserId, &user.UserName, &user.Email); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil

}

func (q *Query) CheckUserIdExists(userId string) (bool, error) {
	query := `SELECT EXISTS (
				SELECT 1 FROM users WHERE user_id=$1
			 )`

	var isUserIdExists bool

	if err := q.db.QueryRow(query, userId).Scan(&isUserIdExists); err != nil {
		return false, err
	}

	return isUserIdExists, nil
}

func (q *Query) UpdateNewPassword(userId string, password string) error {
	query := `UPDATE users SET password = $2 WHERE user_id = $1`
	_, err := q.db.Exec(query, userId, password)
	return err
}

func (q *Query) CheckUserEmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(
				SELECT 1 FROM users WHERE email=$1
			)`
	var isEmailExists bool

	if err := q.db.QueryRow(query, email).Scan(&isEmailExists); err != nil {
		return false, err
	}
	return isEmailExists, nil
}

func (q *Query) StoreOtp(email string, otp string) error {
	query := `INSERT into otps(email,otp) VALUES ($1,$2) ON CONFLICT (email) DO UPDATE 
			  SET otp = EXCLUDED.otp`

	_, err := q.db.Exec(query, email, otp)

	return err
}

func (q *Query) ClearOtp(email string, otp string) error {
	query := `DELETE from otps WHERE email=$1 AND otp=$2`
	_, err := q.db.Exec(query, email, otp)
	return err
}

func (q *Query) IsOtpValid(email string, otp string) (bool, string, error) {
	query1 := `SELECT EXISTS(
				SELECT 1 FROM otps WHERE email = $1 AND otp = $2
			)`
	query2 := `DELETE FROM otps WHERE email = $1 AND otp = $2`

	query3 := `SELECT user_id FROM users WHERE email = $1`

	tx, err := q.db.Begin()

	if err != nil {
		return false, "", err
	}

	var isRowExists bool

	if err := tx.QueryRow(query1, email, otp).Scan(&isRowExists); err != nil {
		tx.Rollback()
		return false, "", err
	}

	if _, err := tx.Exec(query2, email, otp); err != nil {
		tx.Rollback()
		return false, "", err
	}

	var userId string

	if err := tx.QueryRow(query3, email).Scan(&userId); err != nil {
		tx.Rollback()
		return false, "", err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return false, "", err
	}

	return isRowExists, userId, nil

}

func (q *Query) GetBiometricDevicesForRegisterForm(userId string) ([]string, error) {
	query := `SELECT unit_id FROM biometric WHERE user_id=$1`

	rows, err := q.db.Query(query, userId)

	if err != nil {
		return nil, err
	}

	var units []string

	for rows.Next() {
		var unit string
		if err := rows.Scan(&unit); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, nil
}

func (q *Query) GetStudentUnitIdsForRegisterForm(unitId string) ([]string, error) {

	var studentUnitIds []string

	query := fmt.Sprintf(`SELECT student_unit_id FROM %s`, unitId)

	rows, err := q.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var studentUnitId string
		if err := rows.Scan(&studentUnitId); err != nil {
			return nil, err
		}

		studentUnitIds = append(studentUnitIds, studentUnitId)
	}

	return studentUnitIds, nil
}
