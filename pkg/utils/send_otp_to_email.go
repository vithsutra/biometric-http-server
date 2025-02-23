package utils

import (
	"errors"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendOtpToEmail(email string, otp string, expireTime string) error {
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

	htmlTemplate := GetForgotPasswordEmailTemplate(otp, expireTime)

	message := gomail.NewMessage()

	message.SetHeader("From", hostEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Your Secure OTP Code is Ready!")
	message.SetBody("text/html", htmlTemplate)

	if err := client.DialAndSend(message); err != nil {
		return err
	}

	return nil

}
