package main

import (
	"fmt"
	"strconv"
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

func (simpleMovingAverage *SimpleMovingAverage) getSimpleMovingAverageTimestamps() error {
	fmt.Println("Getting Stock Prices for simple moving average")
	statement := fmt.Sprintf(`SELECT 
															ticker_prices.id, 
														  ticker_prices.stock_timestamp, 
															ticker_prices.close, 
															ticker_prices.high, 
															ticker_prices.low, 
															ticker_prices.open, 
															ticker_prices.volume, 
															tickers.id, 
															tickers.symbol, 
															tickers.name, 
															tickers.added_timestamp, 
															tickers.updated_timestamp, 
															exchange.id, 
															exchange.name, 
															to_timestamp(stock_timestamp) 
													  FROM 
															ticker_prices, 
															tickers, 
															exchange 
														WHERE 
														  ticker_prices.ticker_id = tickers.id 
														AND 
															tickers.exchange_id = exchange.id 
														AND 
															ticker_id = 1 
														AND 
														  stock_timestamp <= 1469639940 
														ORDER BY 
														  stock_timestamp 
														DESC 
														limit 500;`, tickers.Id)
	//TODO: Error Checking
	prices, _, _ := postgresql_access.QueryDatabase(db, statement)
	price_list := []Prices{}

	for _, price := range prices {
		single_price := Prices{}
		single_price.Id, _ = strconv.Atoi(price[0].(string))
		single_price.Timestamp, _ = strconv.Atoi(price[1].(string))
		single_price.Close, _ = strconv.ParseFloat(price[2].(string), 64)
		single_price.High, _ = strconv.ParseFloat(price[3].(string), 64)
		single_price.Low, _ = strconv.ParseFloat(price[4].(string), 64)
		single_price.Open, _ = strconv.ParseFloat(price[5].(string), 64)
		single_price.Volume, _ = strconv.Atoi(price[6].(string))

		single_price.Ticker.Id, _ = strconv.Atoi(price[7].(string))
		single_price.Ticker.Symbol = price[8].(string)
		single_price.Ticker.Name = price[9].(string)
		single_price.Ticker.AddedTimestamp, _ = strconv.Atoi(price[10].(string))
		single_price.Ticker.Timestamp, _ = strconv.Atoi(price[11].(string))

		single_price.Ticker.Exchange.Id, _ = strconv.Atoi(price[12].(string))
		single_price.Ticker.Exchange.Name = price[13].(string)

		price_list = append(price_list, single_price)
	}
	fmt.Println("Returning Stock prices ->", len(price_list))
	return price_list
}
