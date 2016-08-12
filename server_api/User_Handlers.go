package main

import (
	"fmt"
	"net/http"
)

func startup(w http.ResponseWriter, r *http.Request) {
	//TODO: Only allow this to run if there is no content available. gives full access
	fmt.Fprintln(w, "startup Not Implimented")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User Not Created Not Implimented")
}

func createRole(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Role Not Created Not Implimented")
}

func createPermission(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Permissions Not Created Not Implimented")
}

func genkey(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}
	userKeys.generateKey()
}
