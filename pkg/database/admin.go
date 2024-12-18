package database

import "github.com/VsenseTechnologies/biometric_http_server/internals/models"



func (q *Query) GiveUserAccess(userId string) (models.Admin , error) {
	var user models.Admin
	if err := q.db.QueryRow("SELECT user_name,password FROM users WHERE user_id=$1", userId).Scan(&user.UserName , &user.Password); err != nil {
		return user,err
	}
	return user , nil
}

func (q *Query) FetchAllUsers() ([]models.Admin , error) {
	res , err := q.db.Query("SELECT user_id , user_name FROM users")
	if err != nil {
		return nil,err
	}
	defer res.Close()
	var user models.Admin
	var users []models.Admin
	for res.Next(){
		if err := res.Scan(&user.UserId , &user.UserName); err != nil {
			return nil , err
		}
		users = append(users, user)
	}
	if res.Err() != nil {
		return nil , err
	}
	return users , nil
}