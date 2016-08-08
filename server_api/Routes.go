package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = AccessLog(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"test",
		"POST",
		"/test",
		test,
	},
	Route{
		"create tables",
		"POST",
		"/api/admin/createtables",
		createTables,
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
