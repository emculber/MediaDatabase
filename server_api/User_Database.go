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

func (role *Role) RegisterNewRole() error {
	err := db.QueryRow(`insert into role (role) values ($1) returning id`, role.Role).Scan(&role.Id)
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
