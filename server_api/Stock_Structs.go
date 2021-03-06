package main

import "errors"

type Exchange struct {
	Id   int
	Name string
}

type Tickers struct {
	Id             int
	Symbol         string
	Name           string
	Exchange       Exchange
	AddedTimestamp int
	Timestamp      int
}

type Prices struct {
	Id        int
	Ticker    Tickers
	Timestamp int
	Close     float64
	High      float64
	Low       float64
	Open      float64
	Volume    int
}

type SimpleMovingAverage struct {
	Period    int64
	Timestamp int64
	Prices    []Prices
}

type MultiplePrices struct {
	Prices []Prices
	count  int
}

func (exchange *Exchange) OK() error {
	if len(exchange.Name) == 0 {
		return errors.New("No Exchange")
	}
	return nil
}

func (tickers *Tickers) OK() error {
	if len(tickers.Symbol) == 0 {
		return errors.New("No Symbol")
	}
	return nil
}
