package main

import (
	"errors"
	"fmt"
)

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

	fmt.Println(read, write, execute)

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
