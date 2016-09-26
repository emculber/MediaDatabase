package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func retriveSimpleMovingAverageTimestampDayCount(w http.ResponseWriter, r *http.Request) {
	ticker := Tickers{}
	err := json.NewDecoder(r.Body).Decode(&ticker)
	if err != nil {
		http.Error(w, err.Error(), 400)
		fmt.Println(err)
		return
	}
	fmt.Println("Ticker ->", ticker)
	count, _ := ticker.retriveSimpleMovingAverageTimestampDayCountFromDatabase()
	fmt.Println("Current Count of Simple Moving Average Day Timestamps ->", count)

	if err := json.NewEncoder(w).Encode(count); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Simple Moving Average Day Timestamp Count")
		return
	}
}

func getMaximumSimpleMovingAverageTimestamp(w http.ResponseWriter, r *http.Request) {
	ticker := Tickers{}
	err := json.NewDecoder(r.Body).Decode(&ticker)
	if err != nil {
		http.Error(w, err.Error(), 400)
		fmt.Println(err)
		return
	}
	fmt.Println("Ticker ->", ticker)
	maximumTimestamp, _ := ticker.retriveSimpleMovingAverageMaxTimestamp()

	if err := json.NewEncoder(w).Encode(maximumTimestamp); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Wallet")
		return
	}
}

func getSimpleMovingAverageTimestamps(w http.ResponseWriter, r *http.Request) {
	simpleMovingAverageOptions := SimpleMovingAverageOptions{}
	err := json.NewDecoder(r.Body).Decode(&simpleMovingAverageOptions)
	if err != nil {
		http.Error(w, err.Error(), 400)
		fmt.Println(err)
		return
	}
	return
	maximumTimestamp, _ := ticker.retriveSimpleMovingAverageMaxTimestamp()

	if err := json.NewEncoder(w).Encode(maximumTimestamp); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Wallet")
		return
	}
}
