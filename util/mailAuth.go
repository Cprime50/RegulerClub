package util

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(subject string, to string, html string, name string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", os.Getenv("EMAIL"))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", html)

	dailer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL"), os.Getenv("SMTP_PASSWORD_GMAIL"))

	err := dailer.DialAndSend(msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}
