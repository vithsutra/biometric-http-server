package repository

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)

type AdminRepo struct{
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{
		db,
	}
}

func (ar *AdminRepo) GiveUserAccess(r *http.Request) error {
	var user models.Admin
	var password string
	if err := utils.Decode(r , &user); err != nil {
		return err
	}
	if err := ar.db.QueryRow("SELECT password FROM users WHERE user_name=$1", newUser.UserName).Scan(&password); err != nil {
		return err
	}
	if err := utils.CheckPassword(password , user.Password); err != nil {
		return err
	}
	tmpl, err := template.ParseFiles("../pkg/templates/email.layout.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load email template: %v", err)
	}

	// Create the data for the template
	data := struct {
		Subject  string
		UserName string
		Password string
	}{
		Subject:  "Access Granted to VSENSE Fingerprint Software",
		UserName: newUser.UserName,
		Password: newUser.Password,
	}

	// Generate the email body
	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	// Convert the mail body to bytes
	mailMessage := []byte(fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", data.Subject, emailBody.String()))

	// Set up the SMTP client
	mailDetails := smtp.PlainAuth(
		"", 
		os.Getenv("SMTP_USERNAME"), 
		os.Getenv("SMTP_PASSWORD"), 
		os.Getenv("SMTP_SERVICE"),
	)
	err = smtp.SendMail(
		"smtp.gmail.com:587", 
		mailDetails, 
		os.Getenv("SMTP_USERNAME"), 
		[]string{newUser.Email}, 
		mailMessage,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}