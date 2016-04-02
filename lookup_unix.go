// Not tested on all of these, but since it written so similarly to os/user
// it seems likely that these should all work.
//
// +build darwin dragonfly freebsd !android,linux netbsd openbsd solaris
// +build cgo

package group

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

// #include <unistd.h>
// #include <sys/types.h>
// #include <grp.h>
// #include <stdlib.h>
//
// char* getmember(struct group* grp, int i) {
//   return grp->gr_mem[i];
// }
//
// static int mygetgrgid_r(int gid, struct group *grp,
//   char *buf, size_t buflen, struct group **result) {
//   return getgrgid_r(gid, grp, buf, buflen, result);
// }
//
// static int mygetgrnam_r(const char *name, struct group *grp,
//   char *buf, size_t buflen, struct group **result) {
//   return getgrnam_r(name, grp, buf, buflen, result);
// }
import "C"

func current() (*Group, error) {
	return lookupUnix(syscall.Getgid(), "", false)
}

func lookup(groupname string) (*Group, error) {
	return lookupUnix(-1, groupname, true)
}

func lookupId(gid string) (*Group, error) {
	i, err := strconv.Atoi(gid)
	if err != nil {
		return nil, err
	}
	return lookupUnix(i, "", false)
}

func lookupUnix(gid int, groupname string, lookupByName bool) (*Group, error) {
	var grp C.struct_group
	var result *C.struct_group

	bufSize := C.sysconf(C._SC_GETGR_R_SIZE_MAX)
	if bufSize == -1 {
		bufSize = 1024
	}
	if bufSize <= 0 || bufSize > 1<<20 {
		return nil, fmt.Errorf("group: unreasonable _SC_GETGR_R_SIZE_MAX of %d", bufSize)
	}
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)

	var rv C.int
	if lookupByName {
		nameC := C.CString(groupname)
		defer C.free(unsafe.Pointer(nameC))
		// mygetgrnam_r is a wrapper around getgrnam_r to avoid using size_t, as
		// the Go standard library mentions that C.size_t(bufSize) doesn't work on
		// Solaris.
		rv = C.mygetgrnam_r(nameC, &grp, (*C.char)(buf), C.size_t(bufSize), &result)
		if rv != 0 {
			return nil, fmt.Errorf("group: lookup groupname %s: %s", groupname, syscall.Errno(rv))
		}
		if result == nil {
			return nil, UnknownGroupError(groupname)
		}
	} else {
		// mygetgrgid_r is a wrapper around getgrgid_r to avoid using gid_t, as
		// C.gid_t(gid) for unknown reasons doesn't work on linux.
		rv = C.mygetgrgid_r(C.int(gid), &grp, (*C.char)(buf), C.size_t(bufSize), &result)
		if rv != 0 {
			return nil, fmt.Errorf("group: lookup groupid %d: %s", gid, syscall.Errno(rv))
		}
		if result == nil {
			return nil, UnknownGroupIdError(gid)
		}
	}

	var members []string
	for i := 0; ; i++ {
		// getmember accesses members of a char**, because no mechanism was made
		// available by Cgo for doing so.
		member := C.getmember(&grp, C.int(i))
		if member == nil {
			break
		}
		members = append(members, C.GoString(member))
	}

	return &Group{
		Gid:     strconv.Itoa(int(grp.gr_gid)),
		Name:    C.GoString(grp.gr_name),
		Members: members,
	}, nil
}
