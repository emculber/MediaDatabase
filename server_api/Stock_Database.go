package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/emculber/database_access/postgresql"
)

var stockDatabaseSchema = []string{
	"CREATE TABLE exchange(id SERIAL PRIMARY KEY, name VARCHAR(60))",
	"CREATE TABLE tickers(id SERIAL PRIMARY KEY, symbol VARCHAR(10), name VARCHAR(256), exchange_id INTEGER REFERENCES exchange(id))",
}

var stockDropDatabaseSchema = []string{
	"DROP TABLE tickers",
	"DROP TABLE exchange",
}

func CreateStockTables() {
	//TODO: check if table exists
	for _, table := range stockDatabaseSchema {
		log.WithFields(log.Fields{
			"Table": table,
		}).Info("Creating Table")
		err := postgresql_access.CreateDatabaseTable(db, table)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Creating Table")
		}
	}
}

func DropStockTables() {
	//TODO: check if table exists
	for _, table := range stockDropDatabaseSchema {
		log.WithFields(log.Fields{
			"Table": table,
		}).Info("Drop Table")
		err := postgresql_access.CreateDatabaseTable(db, table)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Drop Table")
		}
	}
}

func (exchange *Exchange) RegisterNewExchange() error {
	err := db.QueryRow(`insert into exchange (name) values($1) returning id`, exchange.Exchange).Scan(&exchange.Id)
	if err != nil {
		return err
	}
	return nil
}

func (tickers *Tickers) RegisterNewTicker() error {
	err := db.QueryRow(`insert into tickers (symbol, name, exchange_id) values($1, $2, $3) returning id`, tickers.Symbol, tickers.Name, tickers.Exchange.Id).Scan(&tickers.Id)
	if err != nil {
		return err
	}
	return nil
}
