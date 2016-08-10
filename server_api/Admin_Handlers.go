package main

import (
	"fmt"
	"net/http"
)

func createTables(w http.ResponseWriter, r *http.Request) {
	//TODO: Register tables or functions to be create to spereate out dependencies
	//TODO: Check for existing tables
	//TODO: Only allow of superadmin/one person
	fmt.Println("Creating Tables")
	CreateTables()
	CreateFinanceTables()
}
