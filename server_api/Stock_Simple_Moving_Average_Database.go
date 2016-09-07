package main

import "fmt"

func (ticker *Tickers) retriveSimpleMovingAverageMaxTimestamp() (int, error) {
	fmt.Println("Getting Simple Moving Average Max Timestamp For Ticker ->", ticker)

	var maximumTimestamp int
	err := db.QueryRow(`SELECT max(sma_timestamp) FROM simple_moving_average WHERE ticker_id=1`, ticker.Id).Scan(&maximumTimestamp)

	if err != nil {
		return 0, err
	}
	return maximumTimestamp, nil
}
