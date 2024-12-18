package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

func (q *Query) Register(user models.Auth , usertype string) error {

	tx , err := q.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}else{
			tx.Commit()
		}
	} ()

	switch(usertype){
	case "admin":
		_ , err = tx.Exec("INSERT INTO admin(user_id , user_name , password) VALUES($1 , $2 , $3)" , user.UserId , user.Name , user.Password)
		if err != nil {
			return err
		}
	case "users":
		_ , err = tx.Exec("INSERT INTO users(user_id , user_name , password) VALUES($1 , $2 , $3)" , user.UserId , user.Name , user.Password)
		_ , err = tx.Exec("INSERT INTO times(user_id,morning_start , morning_end , afternoon_start , afternoon_end , evening_start , evening_end) VALUES($1,$2,$2,$2,$2,$2,$2)" , user.UserId , "00:00:00")
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid user type")
	}
	return nil
}

func (q *Query) Login(requser models.Auth, usertype string) (models.Auth, error) {
	var user models.Auth
	switch usertype {
	case "admin":
		err := q.db.QueryRow(
			"SELECT user_id, user_name, password FROM admin WHERE user_name = $1",
			requser.Name,
		).Scan(&user.UserId, &user.Name, &user.Password)
		if err != nil {
			return user, fmt.Errorf("admin login failed: %w", err)
		}
	case "users":
		err := q.db.QueryRow(
			"SELECT user_id, user_name, password FROM users WHERE user_name = $1",
			requser.Name,
		).Scan(&user.UserId, &user.Name, &user.Password)
		if err != nil {
			return user, fmt.Errorf("user login failed: %w", err)
		}
	default:
		return user, fmt.Errorf("invalid user type: %s", usertype)
	}
	return user, nil
}