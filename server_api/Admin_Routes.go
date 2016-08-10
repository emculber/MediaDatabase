package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type AdminRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type AdminRoutes []AdminRoute

func NewAdminRouter() *mux.Router {

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

var adminRoutes = AdminRoutes{
	AdminRoute{
		"create tables",
		"POST",
		"/api/admin/createtables",
		createTables,
	},
}

