package main

import "errors"

type Role struct {
	Id   int
	Role string
}

type Permissions struct {
	Id         int
	Permission string
}

func (role *Role) OK() error {
	if len(role.Role) == 0 {
		return errors.New("No Role")
	}
	return nil
}

func (permissions *Permissions) OK() error {
	if len(permissions.Permission) == 0 {
		return errors.New("No Permission")
	}
	return nil
}
