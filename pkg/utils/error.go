package utils

import (
	log "github.com/sirupsen/logrus"
)

func FatalWhenError(err error) {
	if err != nil {
		log.WithError(err).Fatal()
	}
}
