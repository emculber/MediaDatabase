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

type User struct {
	Id       int
	Username string
}

type RolePermissions struct {
	Id         int
	Role       Role
	Permission Permissions
	access     int
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

func (rolePermissions *RolePermissions) OK() error {
	if len(rolePermissions.Role.Role) == 0 {
		return errors.New("No Role")
	}
	if len(rolePermissions.Permission.Permission) == 0 {
		return errors.New("No Permission")
	}
	return nil
}

func (user *User) OK() error {
	if len(user.Username) == 0 {
		return errors.New("No Username")
	}
	return nil
}
