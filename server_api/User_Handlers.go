package main

import (
	"fmt"
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User Not Created Not Implimented")
}
