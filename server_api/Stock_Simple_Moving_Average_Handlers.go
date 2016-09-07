package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

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
