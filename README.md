A group-oriented compliment to the os/user package in the Golang standard
library. Currently only written for Unix-compatible operating systems.

Usage is largely similar to the os/user package:

    import "github.com/hegh/group"
    currentGroup, err := group.Current()
    rootById, err := group.LookupId(0)
    rootByName, err := group.Lookup("root")
