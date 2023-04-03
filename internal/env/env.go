package env

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const Project = "democracy-tools"

func GetSmtp() string {

	const key = "SMTP"
	res := os.Getenv(key)
	if res == "" {
		res = "smtp.gmail.com"
	}
	log.Debugf("%s: %s", key, res)

	return res
}

func GetEmailFrom() string {

	return failIfEmpty("EMAIL_FROM")
}

func GetEmailPassword() string {

	return failIfEmpty("EMAIL_PASSWORD")
}

func failIfEmpty(key string) string {

	res := os.Getenv(key)
	if res == "" {
		log.Fatal("Please, add environment variable '%s'", key)
	}
	log.Debugf("%s: %s", key, res)

	return res
}
