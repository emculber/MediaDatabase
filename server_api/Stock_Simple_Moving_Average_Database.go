package main

import (
	"fmt"
	"time"
)

func (ticker *Tickers) retriveSimpleMovingAverageTimestampDayCountFromDatabase() (int, error) {
	fmt.Println("Getting Day Timestamp count")

	currentDate := time.Unix(int64(ticker.Timestamp), 0)
	lowerDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
	upperDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 59, 0, 0, currentDate.Location())

	var amountTimestamp int
	err := db.QueryRow(`SELECT count(id) FROM simple_moving_average WHERE ticker_id=$1 AND simple_moving_average.sma_timestamp > $2 AND simple_moving_average.sma_timestamp < $3`, ticker.Id, lowerDate.Unix(), upperDate.Unix()).Scan(&amountTimestamp)

	if err != nil {
		return 0, err
	}
	return amountTimestamp, nil
}

func (ticker *Tickers) retriveSimpleMovingAverageMaxTimestamp() (int, error) {
	fmt.Println("Getting Simple Moving Average Max Timestamp For Ticker ->", ticker)

	var maximumTimestamp int
	err := db.QueryRow(`SELECT max(sma_timestamp) FROM simple_moving_average WHERE ticker_id=1`, ticker.Id).Scan(&maximumTimestamp)

	if err != nil {
		return 0, err
	}
	return maximumTimestamp, nil
}
