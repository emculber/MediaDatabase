package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/emculber/database_access/postgresql"
)

var userDatabaseSchema = []string{
	"CREATE TABLE registered_user(id serial primary key, username varchar)",
	"CREATE TABLE role(id serial primary key, role varchar)",
	"CREATE TABLE permissions(id serial primary key, permission varchar)",
	"CREATE TABLE role_permissions(id serial primary key, role_id integer references role(id), permissions_id integer references permissions(id), access varchar)",
	"CREATE TABLE user_keys(id serial primary key, user_id integer references registered_user(id), role_permissions_id integer references role_permissions(id), key varchar)",
}

func CreateUserTables() {
	for _, table := range userDatabaseSchema {
		log.WithFields(log.Fields{
			"Table": table,
		}).Info("Creating Table")
		err := postgresql_access.CreateDatabaseTable(db, table)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Creating Table")
		}
	}
}

func (userCounts *UserCounts) GetUserCounts() error {
	err := db.QueryRow(`SELECT count(registered_user.id) FROM registered_user`).Scan(&userCounts.UserCount)
	if err != nil {
		userCounts.UserCount = 0
	}
	err = db.QueryRow(`SELECT count(role.id) FROM role`).Scan(&userCounts.RoleCount)
	if err != nil {
		userCounts.RoleCount = 0
	}
	err = db.QueryRow(`SELECT count(permissions.id) FROM permissions`).Scan(&userCounts.PermissionsCount)
	if err != nil {
		userCounts.PermissionsCount = 0
	}
	err = db.QueryRow(`SELECT count(role_permissions.id) FROM role_permissions`).Scan(&userCounts.RolePermissionsCount)
	if err != nil {
		userCounts.RolePermissionsCount = 0
	}
	err = db.QueryRow(`SELECT count(user_keys.id) FROM user_keys`).Scan(&userCounts.UserKeysCount)
	if err != nil {
		userCounts.UserKeysCount = 0
	}
	return nil
}

func (role *Role) RegisterNewRole() error {
	err := db.QueryRow(`insert into role (role) values ($1) returning id`, role.Role).Scan(&role.Id)
	if err != nil {
		return err
	}
	return nil
}

func (role *Role) GetRoleId() error {
	err := db.QueryRow("SELECT role.id FROM role where role.role=$1", role.Role).Scan(&role.Id)
	if err != nil {
		return err
	}
	return nil
}

func (permissions *Permissions) RegisterNewPermissions() error {
	err := db.QueryRow(`insert into permissions (permission) values ($1) returning id`, permissions.Permission).Scan(&permissions.Id)
	if err != nil {
		return err
	}
	return nil
}

func (permissions *Permissions) GetPermissionsId() error {
	err := db.QueryRow("SELECT permissions.id FROM permissions where permission=$1", permissions.Permission).Scan(&permissions.Id)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) RegisterNewUser() error {
	err := db.QueryRow(`insert into registered_user (username) values ($1) returning id`, user.Username).Scan(&user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (rolePermissions *RolePermissions) RegisterNewRolePermissions() error {
	err := db.QueryRow(`insert into role_permissions (role_id, permissions_id, access) values ($1, $2, $3) returning id`, rolePermissions.Role.Id, rolePermissions.Permission.Id, rolePermissions.access).Scan(&rolePermissions.Id)
	if err != nil {
		return err
	}
	return nil
}

func (userKeys *UserKeys) RegisterNewUserKeys() error {
	err := db.QueryRow(`insert into user_keys (user_id, role_permissions_id, key) values ($1, $2, $3) returning id`, userKeys.User.Id, userKeys.RolePermissions.Id, userKeys.Key).Scan(&userKeys.Id)
	if err != nil {
		return err
	}
	return nil
}
