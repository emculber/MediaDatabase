package main

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func startup(w http.ResponseWriter, r *http.Request) {
	//TODO: Only allow this to run if there is no content available. gives full access
	userCounts := UserCounts{}
	userCounts.GetUserCounts()
	userKeys := UserKeys{}
	log.WithFields(log.Fields{
		"Role Count":             userCounts.RoleCount,
		"Permissions Count":      userCounts.PermissionsCount,
		"User Count":             userCounts.UserCount,
		"Role Permissions Count": userCounts.RolePermissionsCount,
		"User Keys Count":        userCounts.UserKeysCount,
	}).Info("User Counts")
	if userCounts.RoleCount == 0 {
		role := Role{Role: "admin"}
		if err := role.RegisterNewRole(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Registering New Role")
			return
		}
		userKeys.RolePermissions.Role = role
	}
	if userCounts.PermissionsCount == 0 {
		permissions := Permissions{Permission: "all"}
		if err := permissions.RegisterNewPermissions(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Registering New Role")
			return
		}
		userKeys.RolePermissions.Permission = permissions
	}
	if userCounts.UserCount == 0 {
		user := User{Username: "admin"}
		if err := user.RegisterNewUser(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Registering New Role")
			return
		}
		userKeys.User = user
	}
	if userCounts.RolePermissionsCount == 0 {
		userKeys.RolePermissions.access = 7
		if err := userKeys.RolePermissions.RegisterNewRolePermissions(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Registering New Role")
			return
		}
	}
	if userCounts.UserKeysCount == 0 {
		userKeys.generateKey()

		if err := userKeys.RegisterNewUserKeys(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Registering New Role")
			return
		}
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}
	user := User{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")
	user.Username = r.PostFormValue("username")

	if err := user.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error on user OK")
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

	if err := user.RegisterNewUser(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Role")
		return
	}

	w.Write([]byte("OK"))
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

func createAccess(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}
	rolePermissions := RolePermissions{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")
	rolePermissions.Role.Role = r.PostFormValue("role")
	rolePermissions.Permission.Permission = r.PostFormValue("permission")
	rolePermissions.access, _ = strconv.Atoi(r.PostFormValue("access"))

	if err := rolePermissions.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error on role permissions OK")
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

	if err := rolePermissions.Role.GetRoleId(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Role Permissions Role")
		return
	}

	if err := rolePermissions.Permission.GetPermissionsId(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Role Permissions Permissions")
		return
	}

	if err := rolePermissions.RegisterNewRolePermissions(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Role")
		return
	}

	w.Write([]byte("OK"))
}

func createUserKeys(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}
	newUserKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")
	newUserKeys.RolePermissions.Id, _ = strconv.Atoi(r.PostFormValue("role_permissions_id"))
	newUserKeys.User.Id, _ = strconv.Atoi(r.PostFormValue("user_id"))

	if err := userKeys.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error no user keys OK")
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

	newUserKeys.generateKey()

	if err := newUserKeys.RegisterNewUserKeys(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Role")
		return
	}

	w.Write([]byte("OK"))
}
