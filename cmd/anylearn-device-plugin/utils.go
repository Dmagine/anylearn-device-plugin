package main

import log "github.com/sirupsen/logrus"

func fatalWhenError(err error) {
	if err != nil {
		log.WithError(err).Fatal()
	}
}
