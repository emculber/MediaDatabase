package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func createExchanges(w http.ResponseWriter, r *http.Request) {
	exchange := Exchange{}

	exchange.Name = "NASDAQ"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}
	exchange.Name = "NYSE MKT"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}
	exchange.Name = "New York Stock Exchange"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}
	exchange.Name = "NYSE ARCA"
	if err := exchange.RegisterNewExchange(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Exchange")
		return
	}
	exchange.Name = "BATS Global Markets"
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

	err := json.NewDecoder(r.Body).Decode(&ticker)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Println(ticker)
	if err := ticker.RegisterNewTicker(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Ticker")
		return
	}
}

func sendTickers(w http.ResponseWriter, r *http.Request) {
	/*
		tickerList := TickerList{}

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&tickerList)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		tickers := tickerList.Tickers

		currentTickers := getTickers()

		currentTickersLength := len(currentTickers)
		newTickersLength := len(tickers)

		for currentIndex := 0; currentIndex < currentTickersLength; currentIndex++ {
			for newIndex := 0; newIndex < newTickersLength; newIndex++ {
				if currentTickers[currentIndex].Symbol == tickers[newIndex].Symbol &&
					currentTickers[currentIndex].Name == tickers[newIndex].Name &&
					currentTickers[currentIndex].Exchange.Id == tickers[newIndex].Exchange.Id {

					currentTickers = append(currentTickers[:currentIndex], currentTickers[currentIndex+1:]...)
					currentTickersLength--
					currentIndex--

					tickers = append(tickers[:newIndex], tickers[newIndex+1:]...)
					newTickersLength--
					newIndex--

					/*
						fmt.Println("Update Current Ticker Length ->", currentTickersLength)
						fmt.Println("Update Current Index ->", currentIndex)
						fmt.Println("Update New Ticker Length ->", newTickersLength)
						fmt.Println("Update New Index ->", newIndex)
					break
				}

			}
		}

		fmt.Println("Update Current Ticker Length ->", currentTickersLength)
		fmt.Println("Update New Ticker Length ->", newTickersLength)

		//TODO: Impliment reverse archive

		tickerAudit := TickerAudit{
			AuditTimestamp:            time.Now().Unix(),
			TickerListUpdateTimestamp: tickerList.Timestamp,
			AddedCount:                newTickersLength,
			ChangeCount:               currentTickersLength,
		}
		if err := tickerAudit.RegisterNewAudit(); err != nil {
			log.WithFields(log.Fields{
				"Ticker Audit": tickerAudit,
				"Error":        err,
			}).Error("Error Registering New Ticker Audit")
			return
		}

		for _, changeTicker := range currentTickers {
			changeTicker.Archived = "Y"

			if err := changeTicker.RegisterNewTicker(); err != nil {
				log.WithFields(log.Fields{
					"Ticker": changeTicker,
					"Error":  err,
				}).Error("Error Registering Change Ticker")
			}
			tickerUpdate := TickerUpdate{
				TickerAudit:     tickerAudit,
				UpdateTimestamp: time.Now().Unix(),
				UpdateType:      "Archived",
				Ticker:          changeTicker,
			}

			if err := tickerUpdate.RegisterNewTickerUpdate(); err != nil {
				log.WithFields(log.Fields{
					"Ticker Update": tickerUpdate,
					"Error":         err,
				}).Error("Error Registering New Ticker Update")
			}
		}
		for _, newTicker := range tickers {

			newTicker.Archived = "N"

			if err := newTicker.RegisterNewTicker(); err != nil {
				log.WithFields(log.Fields{
					"Ticker": newTicker,
					"Error":  err,
				}).Error("Error Registering New Ticker")
			}

			tickerUpdate := TickerUpdate{
				TickerAudit:     tickerAudit,
				UpdateTimestamp: time.Now().Unix(),
				UpdateType:      "Add",
				Ticker:          newTicker,
			}

			if err := tickerUpdate.RegisterNewTickerUpdate(); err != nil {
				log.WithFields(log.Fields{
					"Ticker Update": tickerUpdate,
					"Error":         err,
				}).Error("Error Registering New Ticker Update")
			}
		}

		w.Write([]byte("OK"))
	*/
}

func getAllTickers(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(getTickers()); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Wallet")
		return
	}
}

func createPrices(w http.ResponseWriter, r *http.Request) {
	prices := []Prices{}

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&prices)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	for i, price := range prices {
		if i == 0 {
			log.WithFields(log.Fields{
				"Ticker": prices[0].Ticker,
			}).Info("Registering New Prices With Ticker")
		}
		if err := price.RegisterNewPrice(); err != nil {
			log.WithFields(log.Fields{
				"Price": price,
				"Index": i,
				"Error": err,
			}).Error("Error Registering New Price")
		}
	}

	w.Write([]byte("OK"))
}
