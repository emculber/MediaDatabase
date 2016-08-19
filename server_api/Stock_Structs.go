package main

import "errors"

type Exchange struct {
	Id       int
	Exchange string
}

type Tickers struct {
	Id       int
	Symbol   string
	Name     string
	Exchange Exchange
}

func (exchange *Exchange) OK() error {
	if len(exchange.Exchange) == 0 {
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
