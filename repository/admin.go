package repository

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"text/template"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
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

func (ar *AdminRepo) FetchAllUsers() ([]models.Admin,error) {
	query := database.NewQuery(ar.db)
	users , err := query.FetchAllUsers()
	if err != nil {
		return nil , err
	}
	return users , nil
}
func (ar *AdminRepo) GiveUserAccess(r *http.Request) error {
	var user models.Admin
	if err := utils.Decode(r, &user); err != nil {
		return err
	}
	query := database.NewQuery(ar.db)
	dbusr, err := query.GiveUserAccess(user.UserId)
	if err != nil {
		return err
	}
	if err := utils.CheckPassword(dbusr.Password, user.Password); err != nil {
		return err
	}
	go func(user models.Admin, dbusr models.Admin) {
		tmpl, err := template.ParseFiles("pkg/templates/email.layout.tmpl")
		if err != nil {
			log.Printf("Unable to parse email template: %v\n" , err)
			return
		}

		data := struct {
			Subject  string
			UserName string
			Password string
		}{
			Subject:  "Access Granted to VSENSE Fingerprint Software",
			UserName: dbusr.UserName,
			Password: user.Password,
		}

		var emailBody bytes.Buffer
		if err := tmpl.Execute(&emailBody, data); err != nil {
			log.Printf("Error executing email template: %v\n", err)
			return
		}

		mailMessage := []byte(fmt.Sprintf(
			"Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
			data.Subject, emailBody.String(),
		))

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
			[]string{user.Email},
			mailMessage,
		)
		if err != nil {
			log.Printf("Error sending email: %v\n", err)
		}
	}(user, dbusr)
	return nil
}
