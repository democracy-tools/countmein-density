package internal

import log "github.com/sirupsen/logrus"

func validateToken(token string) bool {

	count := len(token)
	if count < 30 || count > 40 {
		log.Infof("invalid token length '%d'", count)
		return false
	}

	return true
}
