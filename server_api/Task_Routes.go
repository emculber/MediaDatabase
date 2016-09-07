package main

import (
	"net/http"
)

type TaskRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type TaskRoutes []TaskRoute

func (router *MuxRouter) TaskRouter() {

	for _, route := range taskRoutes {
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

var taskRoutes = TaskRoutes{
	TaskRoute{
		"New Task",
		"POST",
		"/api/task/new",
		newTask,
	},
	TaskRoute{
		"Get All Tasks",
		"POST",
		"/api/task/all",
		getTasks,
	},
}
