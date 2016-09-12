package main

import (
	"net/http"
)

type StockRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type StockRoutes []StockRoute

func (router *MuxRouter) StockRouter() {

	for _, route := range stockRoutes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = AccessLog(handler, route.Name)

		router.
			Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}

var stockRoutes = StockRoutes{
	StockRoute{
		"Add Exchange",
		"POST",
		"/api/stock/createexchange",
		createExchanges,
	},
	StockRoute{
		"Add ticker",
		"POST",
		"/api/stock/createticker",
		createTicker,
	},
	StockRoute{
		"Add tickers",
		"POST",
		"/api/stock/createtickers",
		sendTickers,
	},
	StockRoute{
		"Get All tickers",
		"GET",
		"/api/stock/gettickers",
		getAllTickers,
	},
	StockRoute{
		"Add prices",
		"POST",
		"/api/stock/createprices",
		createPrices,
	},
	StockRoute{
		"Get All Prices for Ticker",
		"POST",
		"/api/stock/gettickerprices",
		getTickerPrices,
	},
	StockRoute{
		"Get Max Timestamp",
		"GET",
		"/api/stock/getmaxtimestamp",
		getMaxTimestamp,
	},
	StockRoute{
		"Get Min Timestamp",
		"GET",
		"/api/stock/getmintimestamp",
		getMinTimestamp,
	},
	StockRoute{
		"Get Max Timestamp",
		"POST",
		"/api/stock/gettimestampdaycount",
		getTimestampDayCount,
	},
	StockRoute{
		"Get price Timestamps",
		"POST",
		"/api/stock/gettimestampsday",
		getTimestampsDay,
	},
	StockRoute{
		"Get Simple Moving Average Day Timestamp Count",
		"POST",
		"/api/stock/retrivesimplemovingaveragetimestampdaycount",
		retriveSimpleMovingAverageTimestampDayCount,
	},
}
