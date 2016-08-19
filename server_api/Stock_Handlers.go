package main

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func createExchanges(w http.ResponseWriter, r *http.Request) {
	exchange := Exchange{}

	exchange.Exchange = "NASDAQ"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}
	exchange.Exchange = "NYSE MKT"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}
	exchange.Exchange = "New York Stock Exchange"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}
	exchange.Exchange = "NYSE ARCA"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}
	exchange.Exchange = "BATS Global Markets"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}

	w.Write([]byte("OK"))
}

func createTicker(w http.ResponseWriter, r *http.Request) {
	ticker := Tickers{}

	r.ParseForm()
	ticker.Symbol = r.PostFormValue("symbol")
	ticker.Name = r.PostFormValue("name")
	ticker.Exchange.Id, _ = strconv.Atoi(r.PostFormValue("exchange_id"))

	if err := ticker.RegisterNewTicker(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Ticker")
		return
	}
	w.Write([]byte("OK"))
}
