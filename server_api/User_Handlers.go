package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func startup(w http.ResponseWriter, r *http.Request) {
	//TODO: Only allow this to run if there is no content available. gives full access
	fmt.Fprintln(w, "startup Not Implimented")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User Not Created Not Implimented")
}

func createRole(w http.ResponseWriter, r *http.Request) {
	role := Role{}
	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")
	role.Role = r.PostFormValue("role")

	if err := role.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error on role OK")
		return
	}

	if err := userKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	}

	if err := userKeys.RolePermissions.checkPermissions("user"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	if err := userKeys.RolePermissions.checkAccess("write"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	if err := role.RegisterNewRole(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Role")
		return
	}

	w.Write([]byte("OK"))
}

func createPermission(w http.ResponseWriter, r *http.Request) {
	permissions := Permissions{}
	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")
	permissions.Permission = r.PostFormValue("permission")

	if err := permissions.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error on permissions OK")
		return
	}

	if err := userKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	}

	if err := userKeys.RolePermissions.checkPermissions("user"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	if err := userKeys.RolePermissions.checkAccess("write"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	if err := permissions.RegisterNewPermissions(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Role")
		return
	}

	w.Write([]byte("OK"))
}

func genkey(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}
	userKeys.generateKey()
}
