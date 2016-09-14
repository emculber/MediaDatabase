package main

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emculber/database_access/postgresql"
	"github.com/lib/pq"
)

var stockDatabaseSchema = []string{
	"CREATE TABLE exchange(id SERIAL PRIMARY KEY, name VARCHAR(60))",
	"CREATE TABLE tickers(id SERIAL PRIMARY KEY, symbol VARCHAR(10), name VARCHAR(256), exchange_id INTEGER REFERENCES exchange(id), added_timestamp BIGINT, updated_timestamp BIGINT, UNIQUE (symbol, name, exchange_id))",
	"CREATE TABLE ticker_prices (id SERIAL PRIMARY KEY, ticker_id INTEGER REFERENCES tickers(id), stock_timestamp INTEGER, close REAL, high REAL, low REAL, open REAL, volume INTEGER, UNIQUE(ticker_id, stock_timestamp))",
	"CREATE TABLE simple_moving_average (id SERIAL PRIMARY KEY, ticker_id INTEGER REFERENCES tickers(id), sma_timestamp INTEGER, sma REAL, period INTEGER, UNIQUE(ticker_id, sma_timestamp))",
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
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (tickers *Tickers) getTickerPrices() []Prices {
	fmt.Println("Getting Stock Prices")
	statement := fmt.Sprintf("SELECT ticker_prices.id, ticker_prices.stock_timestamp, ticker_prices.close, ticker_prices.high, ticker_prices.low, ticker_prices.open, ticker_prices.volume, tickers.id, tickers.symbol, tickers.name, tickers.added_timestamp, tickers.updated_timestamp, exchange.id, exchange.name FROM ticker_prices, tickers, exchange WHERE ticker_prices.ticker_id = tickers.id AND tickers.exchange_id = exchange.id AND tickers.id = %d", tickers.Id)
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

func (tickers *Tickers) RegisterNewTicker() error {
	//err := db.QueryRow(`insert into tickers (symbol, name, exchange_id, added_timestamp, updated_timestamp) values($1, $2, $3, $4, $5) returning id`, tickers.Symbol, tickers.Name, tickers.Exchange.Id, tickers.Timestamp, tickers.Timestamp).Scan(&tickers.Id)
	fmt.Println("Setting up columns")
	columns := []string{"symbol", "name", "exchange_id", "added_timestamp", "updated_timestamp"}

	fmt.Println("Converting data to an interfaces")
	data := make([]interface{}, 5)
	data[0] = tickers.Symbol
	data[1] = tickers.Name
	data[2] = tickers.Exchange.Id
	data[3] = tickers.Timestamp
	data[4] = tickers.Timestamp

	fmt.Println("Trying to insert Single Data Values")

	if err := postgresql_access.InsertSingleDataValue(db, "tickers", columns, data); err != nil {
		if err, ok := err.(*pq.Error); ok {
			fmt.Println(err)
			if err.Code == "23505" {
				fmt.Println("Updating Ticker ->", tickers)
				//err := db.QueryRow(`UPDATE tickers SET symbol = $1, name=$2, exchange_id = $3, updated_timestamp = $4 WHERE id = $5 returning id`, tickers.Symbol, tickers.Name, tickers.Exchange.Id, tickers.Timestamp, tickers.Id).Scan(&tickers.Id)
				fmt.Println("Update err ->", err)
				return nil
				if err != nil {
					return err
				}
			} else {
				return err
			}
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

func getTickers() []Tickers {
	fmt.Println("Getting Stock Tickers")
	statement := fmt.Sprintf("SELECT tickers.id, tickers.symbol, tickers.name, exchange.id, exchange.name, tickers.added_timestamp, updated_timestamp from tickers, exchange where tickers.exchange_id = exchange.id")
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
		single_ticker.AddedTimestamp, _ = strconv.Atoi(ticker[5].(string))
		single_ticker.Timestamp, _ = strconv.Atoi(ticker[6].(string))
		ticker_list = append(ticker_list, single_ticker)
	}
	fmt.Println("Returning Stock Tickers ->", len(ticker_list))
	return ticker_list
}

func retriveMaxTimestamp() (int, error) {
	var maxTimestamp int
	err := db.QueryRow(`SELECT max(stock_timestamp) from ticker_prices`).Scan(&maxTimestamp)
	if err != nil {
		return 0, err
	}
	fmt.Println("Getting Max Timestamp ->", maxTimestamp)
	return maxTimestamp, nil
}

func retriveMinTimestamp() (int, error) {
	fmt.Println("Getting Min Timestamp")
	var minTimestamp int
	err := db.QueryRow(`SELECT min(stock_timestamp) from ticker_prices`).Scan(&minTimestamp)
	if err != nil {
		return 0, err
	}
	return minTimestamp, nil
}

func (ticker *Tickers) retriveDayTimestampCount() (int, error) {
	fmt.Println("Getting Day Timestamp count")

	currentDate := time.Unix(int64(ticker.Timestamp), 0)
	lowerDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
	upperDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 59, 0, 0, currentDate.Location())

	var amountTimestamp int
	err := db.QueryRow(`SELECT count(id) FROM ticker_prices WHERE ticker_id=$1 AND stock_timestamp > $2 AND ticker_prices.stock_timestamp < $3`, ticker.Id, lowerDate.Unix(), upperDate.Unix()).Scan(&amountTimestamp)

	if err != nil {
		return 0, err
	}
	return amountTimestamp, nil
}

func (ticker *Tickers) retriveDayTimestamp() []Prices {
	fmt.Println("Getting Day Timestamps")

	currentDate := time.Unix(int64(ticker.Timestamp), 0)
	lowerDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
	upperDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 59, 0, 0, currentDate.Location())

	statement := fmt.Sprintf("SELECT ticker_prices.id, ticker_prices.stock_timestamp, ticker_prices.close, ticker_prices.high, ticker_prices.low, ticker_prices.open, ticker_prices.volume, tickers.id, tickers.symbol, tickers.name, tickers.added_timestamp, tickers.updated_timestamp, exchange.id, exchange.name FROM ticker_prices, tickers, exchange WHERE ticker_prices.ticker_id = tickers.id AND tickers.exchange_id = exchange.id AND ticker_id=%d AND stock_timestamp > %d AND ticker_prices.stock_timestamp < %d", ticker.Id, lowerDate.Unix(), upperDate.Unix())
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

func (ticker *Tickers) retriveSMADayTimestampCount() (int, error) {
	fmt.Println("Getting SMA Timestamp Count")

	currentDate := time.Unix(int64(ticker.Timestamp), 0)
	lowerDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
	upperDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 59, 0, 0, currentDate.Location())

	var amountTimestamp int
	err := db.QueryRow(`SELECT count(id) FROM ticker_prices WHERE ticker_id=$1 AND stock_timestamp > $2 AND ticker_prices.stock_timestamp < $3`, ticker.Id, lowerDate.Unix(), upperDate.Unix()).Scan(&amountTimestamp)

	if err != nil {
		return 0, err
	}
	return amountTimestamp, nil
}
