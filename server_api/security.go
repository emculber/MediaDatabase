package main

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

func (rolePermissions *RolePermissions) checkPermissions(permissionNeeded string) error {

	if rolePermissions.Permission.Permission == "all" || rolePermissions.Permission.Permission == "admin" {
		return nil
	}

	if permissionNeeded != rolePermissions.Permission.Permission {
		return errors.New("Access Denied")
	}
	return nil
}

func (rolePermissions *RolePermissions) checkAccess(accessNeeded string) error {
	binaryRepresentation := fmt.Sprintf("%b", rolePermissions.access)

	var read, write, execute bool
	if len(binaryRepresentation) >= 1 {
		read = binaryRepresentation[0] == '1'
	}
	if len(binaryRepresentation) >= 2 {
		write = binaryRepresentation[1] == '1'
	}
	if len(binaryRepresentation) >= 3 {
		execute = binaryRepresentation[2] == '1'
	}

	log.WithFields(log.Fields{
		"Read":    read,
		"Write":   write,
		"Exicute": execute,
	}).Info("Access")

	if accessNeeded == "read" {
		if read {
			return nil
		} else {
			return errors.New("Access Denied")
		}
	} else if accessNeeded == "write" {
		if write {
			return nil
		} else {
			return errors.New("Access Denied")
		}
	} else if accessNeeded == "execute" {
		if execute {
			return nil
		} else {
			return errors.New("Access Denied")
		}
	} else {
		return errors.New("Access Needed is Not Found")
	}
}
