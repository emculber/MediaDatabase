package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func InitLogger() {
	path := os.Getenv("GOPATH")
	var filePath string = path + "/logs/api_logs/media_database.log"

	log.SetFormatter(&log.TextFormatter{})
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	log_file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{
			"File Path": filePath,
			"Error":     err.Error(),
		}).Error("Error Opening File")
	}

	log.SetOutput(log_file)

	log.WithFields(log.Fields{
		"Log Format": "Text Format",
		"Log level":  "Info",
		"Log Output": log_file,
	}).Info("Format, Level, Output set")
}
