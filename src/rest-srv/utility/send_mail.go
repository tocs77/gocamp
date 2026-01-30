package utility

import (
	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("mailhog", 1025, "", "")
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
