package logging

import (
	log "github.com/sirupsen/logrus"
)

// reuse logger configured with app context throughout the app
var Log *log.Entry

func Init(level string, appVersion string) {
	// set log level
	switch level {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	Log = log.WithFields(log.Fields{
		"app":     "rentals-api",
		"version": appVersion,
	})
}
