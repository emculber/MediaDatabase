package main

import (
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/emculber/database_access/postgresql"
	"github.com/lib/pq"
)

var stockDatabaseSchema = []string{
	"CREATE TABLE exchange(id SERIAL PRIMARY KEY, name VARCHAR(60))",
	"CREATE TABLE tickers(id SERIAL PRIMARY KEY, symbol VARCHAR(10), name VARCHAR(256), exchange_id INTEGER REFERENCES exchange(id), added_timestamp BIGINT, updated_timestamp BIGINT, UNIQUE (symbol, name, exchange_id))",
	"CREATE TABLE ticker_prices (id SERIAL PRIMARY KEY, ticker_id INTEGER REFERENCES tickers(id), stock_timestamp INTEGER, close REAL, high REAL, low REAL, open REAL, volume INTEGER, UNIQUE(ticker_id, stock_timestamp))",
}

var stockDropDatabaseSchema = []string{
	"DROP TABLE ticker_prices",
	"DROP TABLE ticker_audit",
	"DROP TABLE ticker_updates",
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
	err := db.QueryRow(`insert into exchange (name) values($1) returning id`, exchange.Name).Scan(&exchange.Id)
	if err != nil {
		return err
	}
	return nil
}

func (prices *Prices) RegisterNewPrice() error {
	err := db.QueryRow(`insert into ticker_prices (ticker_id, stock_timestamp, close, high, low, open, volume) values($1, $2, $3, $4, $5, $6, $7) returning id`, prices.Ticker.Id, prices.Timestamp, prices.Close, prices.High, prices.Low, prices.Open, prices.Volume).Scan(&prices.Id)
	if err != nil {
		return err
	}
	return nil
}

func (tickers *Tickers) RegisterNewTicker() error {
	err := db.QueryRow(`insert into tickers (symbol, name, exchange_id, added_timestamp, updated_timestamp) values($1, $2, $3, $4, $5) returning id`, tickers.Symbol, tickers.Name, tickers.Exchange.Id, tickers.Timestamp, tickers.Timestamp).Scan(&tickers.Id)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			err := db.QueryRow(`UPDATE tickers SET symbol = $1, name=$2, exchange_id = $3, updated_timestamp = $4 WHERE id = $5 returning id`, tickers.Symbol, tickers.Name, tickers.Exchange.Id, tickers.Timestamp, tickers.Id).Scan(&tickers.Id)
			fmt.Println("Update")
			return nil
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	fmt.Println("Insert")
	return nil
}

func (tickers *Tickers) updateTicker() error {
	err := db.QueryRow(`UPDATE tickers SET symbol = $1, name=$2, exchange_id = $3, archived = $4 WHERE id = $5 returning id`, tickers.Symbol, tickers.Name, tickers.Exchange.Id, tickers.Timestamp, tickers.Id).Scan(&tickers.Id)
	if err != nil {
		return err
	}
	return nil
}

func (tickerAudit *TickerAudit) RegisterNewAudit() error {
	err := db.QueryRow(`insert into ticker_audit (audit_timestamp, ticker_list_update_timestamp, added_count, change_count) values($1, $2, $3, $4) returning id`, tickerAudit.AuditTimestamp, tickerAudit.TickerListUpdateTimestamp, tickerAudit.AddedCount, tickerAudit.ChangeCount).Scan(&tickerAudit.Id)
	if err != nil {
		return err
	}
	return nil
}

func (tickerUpdate *TickerUpdate) RegisterNewTickerUpdate() error {
	err := db.QueryRow(`insert into ticker_updates (ticker_audit_id, update_timestamp, update_type, ticker_id) values($1, $2, $3, $4) returning id`, tickerUpdate.TickerAudit.Id, tickerUpdate.UpdateTimestamp, tickerUpdate.UpdateType, tickerUpdate.Ticker.Id).Scan(&tickerUpdate.Id)
	if err != nil {
		return err
	}
	return nil
}

func getTickers() []Tickers {
	fmt.Println("Getting Stock Tickers")
	statement := fmt.Sprintf("SELECT tickers.id, tickers.symbol, tickers.name, exchange.id, exchange.name, tickers.archived from tickers, exchange where tickers.exchange_id = exchange.id")
	//TODO: Error Checking
	tickers, _, _ := postgresql_access.QueryDatabase(db, statement)
	ticker_list := []Tickers{}

	for _, ticker := range tickers {
		single_ticker := Tickers{}
		single_ticker.Id, _ = strconv.Atoi(ticker[0].(string))
		single_ticker.Symbol = ticker[1].(string)
		single_ticker.Name = ticker[2].(string)
		single_ticker.Exchange.Id, _ = strconv.Atoi(ticker[3].(string))
		single_ticker.Exchange.Name = ticker[4].(string)
		//single_ticker.Timestamp = ticker[5].(string)
		ticker_list = append(ticker_list, single_ticker)
	}
	fmt.Println("Returning Stock Tickers")
	return ticker_list
}
