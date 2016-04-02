// Package group allows group account lookups by name or id.
//
// This is intended to compliment the os/user package, as groups are not
// supported in the default Golang distribution for some unknown reason.
package group

import "strconv"

type Group struct {
	Gid     string // group id
	Name    string
	Members []string
}

type UnknownGroupIdError int

func (e UnknownGroupIdError) Error() string {
	return "group: unknown groupid " + strconv.Itoa(int(e))
}

type UnknownGroupError string

func (e UnknownGroupError) Error() string {
	return "group: unknown group " + string(e)
}
