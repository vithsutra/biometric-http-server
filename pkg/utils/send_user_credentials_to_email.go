package utils

import (
	"errors"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendUserCredentialsToEmail(userName string, password string, email string) error {
	smtpHost := os.Getenv("SMTP_HOST")

	if smtpHost == "" {
		return errors.New("SMTP_HOST env was missing")
	}

	smtpPort := os.Getenv("SMTP_PORT")

	if smtpPort == "" {
		return errors.New("SMTP_PORT env was missing")
	}

	hostEmail := os.Getenv("HOST_EMAIL")

	if hostEmail == "" {
		return errors.New("HOST_EMAIL env was missing")
	}

	appPassword := os.Getenv("APP_PASSWORD")

	if appPassword == "" {
		return errors.New("APP_PASSWORD env was missing")
	}

	to := email

	smtpPortInt, err := strconv.Atoi(smtpPort)

	if err != nil {
		return err
	}

	client := gomail.NewDialer(smtpHost, smtpPortInt, hostEmail, appPassword)

	htmlTemplate := GetUserAccessEmailTemplate(userName, password)

	message := gomail.NewMessage()

	message.SetHeader("From", hostEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Welcome to Vithsutra Technologies")
	message.SetBody("text/html", htmlTemplate)

	if err := client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
