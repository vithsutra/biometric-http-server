package database

import "github.com/VsenseTechnologies/biometric_http_server/internals/models"

func (q *Query) CreateAdmin(admin *models.Admin) error {
	query := `INSERT INTO admin (user_id,user_name,password) VALUES ($1,$2,$3)`
	_, err := q.db.Exec(query, admin.UserId, admin.UserName, admin.Password)
	return err
}

func (q *Query) GetAdminPassword(userName string) (string, error) {
	query := `SELECT password FROM admin WHERE user_name=$1`

	var adminPassword string
	err := q.db.QueryRow(query, userName).Scan(&adminPassword)
	return adminPassword, err
}
