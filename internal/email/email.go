package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"sync"

	"github.com/democracy-tools/countmein-density/internal/env"
	log "github.com/sirupsen/logrus"
)

type singleton struct {
	smtpHost    string
	smtpAddress string
	from        string
	to          string
	password    string
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = newClient()
	})
	return instance
}

func newClient() *singleton {

	host := env.GetSmtp()

	return &singleton{
		smtpHost:    host,
		smtpAddress: fmt.Sprintf("%s:587", host),
		from:        env.GetEmailFrom(),
		to:          env.GetEmailSupport(),
		password:    env.GetEmailPassword(),
	}
}

func (s *singleton) SendError(message string) error {

	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
	body := []byte(fmt.Sprintf("Subject: [CountMeIn] %s\r\n\r\n%s\r\n", message[:70], message))
	return smtp.SendMail(fmt.Sprintf("%s:587", s.smtpHost), auth, s.from, []string{s.to}, body)
}

func (s *singleton) Send(to string, htmlTemplate string, subject string, data any) error {

	body, err := createBody(htmlTemplate, subject, data)
	if err != nil {
		return err
	}

	return smtp.SendMail(s.smtpAddress,
		smtp.PlainAuth("", s.from, s.password, s.smtpHost),
		s.from, []string{to}, body)
}

func createBody(htmlTemplate string, subject string, data any) ([]byte, error) {

	var body bytes.Buffer

	const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subject, mime)))
	t, err := template.ParseFiles(htmlTemplate)
	if err != nil {
		log.Errorf("failed to parse '%s' with '%v'", htmlTemplate, err)
		return nil, err
	}
	err = t.Execute(&body, data)
	if err != nil {
		log.Errorf("failed to execute template '%s' with '%v'", htmlTemplate, err)
		return nil, err
	}

	return body.Bytes(), nil
}
