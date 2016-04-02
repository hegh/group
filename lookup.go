package group

// Current returns the current primary group.
func Current() (*Group, error) {
	return current()
}

// Lookup looks up a group by groupname. If the group cannot be found, the
// returned error is of type UnknownGroupError.
func Lookup(groupname string) (*Group, error) {
	return lookup(groupname)
}

// LookupId looks up a group by groupid. If the group cannot be found, the
// returned error is of type UnknownGroupIdError.
func LookupId(gid string) (*Group, error) {
	return lookupId(gid)
}
