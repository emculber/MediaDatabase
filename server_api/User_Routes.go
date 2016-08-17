package main

import (
	"net/http"
)

type UserRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type UserRoutes []UserRoute

func (router *MuxRouter) UserRouter() {

	for _, route := range userRoutes {
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

var userRoutes = UserRoutes{
	UserRoute{
		"startup",
		"POST",
		"/api/user/startup",
		startup,
	},
	UserRoute{
		"create user",
		"POST",
		"/api/user/createuser",
		createUser,
	},
	UserRoute{
		"create role",
		"POST",
		"/api/user/createrole",
		createRole,
	},
	UserRoute{
		"create permissions",
		"POST",
		"/api/user/createpermission",
		createPermission,
	},
	UserRoute{
		"create access",
		"POST",
		"/api/user/createaccess",
		createAccess,
	},
	UserRoute{
		"create user key",
		"POST",
		"/api/user/createuserkeys",
		createUserKeys,
	},
	UserRoute{
		"Generate Key",
		"POST",
		"/api/user/generatekey",
		genkey,
	},
}
