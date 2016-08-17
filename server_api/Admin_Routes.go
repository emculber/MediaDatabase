package main

import (
	"net/http"
)

type AdminRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type AdminRoutes []AdminRoute

func (router *MuxRouter) AdminRouter() {

	for _, route := range adminRoutes {
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

var adminRoutes = AdminRoutes{
	AdminRoute{
		"create tables",
		"POST",
		"/api/admin/createtables",
		createTables,
	},
	AdminRoute{
		"drop tables",
		"POST",
		"/api/admin/droptables",
		dropTables,
	},
}
