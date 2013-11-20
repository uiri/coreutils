// +build linux openbsd

package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func GetAtimeMtime(name string, noderef bool) (atime, mtime time.Time) {
	fileinfo := new(syscall.Stat_t)
	var err error
	if noderef {
		err = syscall.Lstat(name, fileinfo)
	} else {
		err = syscall.Stat(name, fileinfo)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting info about '%s': %v\n", name, err)
	}
	sec, nsec := fileinfo.Atim.Unix()
	atime = time.Unix(sec, nsec)
	sec, nsec = fileinfo.Mtim.Unix()
	mtime = time.Unix(sec, nsec)
	return atime, mtime
}
