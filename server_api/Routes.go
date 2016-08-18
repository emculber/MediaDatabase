package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (router *MuxRouter) GenericRouter() {

	for _, route := range routes {
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

var routes = Routes{
	Route{
		"test",
		"POST",
		"/test",
		test,
	},
	Route{
		"Add Movie To User Movies",
		"POST",
		"/api/addmovie",
		addMovieToUserMovies,
	},
	Route{
		"Check user key",
		"POST",
		"/api/checkuserkey",
		getUserKey,
	},
	Route{
		"Get All Movies for that User",
		"POST",
		"/api/getmovies",
		getAllMovies,
	},
	Route{
		"Get All",
		"POST",
		"/api/getallregesteredmovies",
		getAllMovies,
	},

	Route{
		"Setup Finance",
		"POST",
		"/api/finance/setup",
		SetupFinance,
	},
	Route{
		"New Finance User",
		"POST",
		"/api/finance/newuser",
		NewFinancalUser,
	},
	Route{
		"New Income",
		"POST",
		"/api/finance/newincome",
		NewIncome,
	},
	Route{
		"Get Income",
		"POST",
		"/api/finance/getincomes",
		GetIncomes,
	},
	Route{
		"New Expense",
		"POST",
		"/api/finance/newexpense",
		NewExpense,
	},
	Route{
		"New Wallet",
		"POST",
		"/api/finance/newwallet",
		NewWallet,
	},
	Route{
		"Get Wallet",
		"POST",
		"/api/finance/getwallets",
		GetWallets,
	},
}
