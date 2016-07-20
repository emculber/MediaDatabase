package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

type ApiLoggerFields struct {
	ip_address  string
	method_type string
}

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

func AccessLog(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.WithFields(log.Fields{
			"Method Type": r.Method,
			"Request Url": r.RequestURI,
			"IP":          r.RemoteAddr,
			"Name":        name,
			"Time":        time.Since(start),
		}).Info("Accessing API")
	})
}
