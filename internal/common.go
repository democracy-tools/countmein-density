package internal

import (
	"regexp"

	"github.com/sirupsen/logrus"
)

const (
	// RegisterUrl = "https://storage.googleapis.com/countmein-web/register.html"
	// VerificationUrl = "https://storage.googleapis.com/countmein-web/phone-verification.html"
	// ObservationUrl  = "https://storage.googleapis.com/countmein-web/demonstration.html"
	JoinUrl = "https://storage.googleapis.com/countmein-web/join-demonstration.html"
)

func validateToken(token string) bool {

	const exp = "^[a-z0-9-]{36}$"
	res, err := regexp.MatchString(exp, token)
	if err != nil {
		logrus.Errorf("failed to compile regexp '%s' with '%v'", exp, err)
		return false
	}
	if !res {
		logrus.Infof("invalid token '%s'", token)
	}

	return res
}
