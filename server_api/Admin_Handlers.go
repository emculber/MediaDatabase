package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func createTables(w http.ResponseWriter, r *http.Request) {
	//TODO: Register tables or functions to be create to spereate out dependencies
	//TODO: Check for existing tables
	//TODO: Only allow of superadmin/one person

	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")

	if err := userKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"User Blame":      userKeys.User.Username,
			"User Permission": userKeys.RolePermissions,
			"User Key":        userKeys.Key,
			"Error":           err,
		}).Error("Error Validating User")
		return
	}

	if err := userKeys.RolePermissions.checkAccess("execute"); err != nil {
		log.WithFields(log.Fields{
			"User Blame":      userKeys.User.Username,
			"User Permission": userKeys.RolePermissions,
			"User Key":        userKeys.Key,
			"Error":           err,
		}).Error("Error Checking Access")
		return
	}

	log.WithFields(log.Fields{
		"User Blame":      userKeys.User.Username,
		"User Permission": userKeys.RolePermissions,
		"User Key":        userKeys.Key,
	}).Info("Running Create Tables")
	CreateTables()
	CreateFinanceTables()
}
