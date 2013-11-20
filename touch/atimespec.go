// +build freebsd netbsd darwin

package main

import (
	"syscall"
	"time"
)

func GetAtimeMtime(name string, noderef bool) (atime, mtime time.Time) {
	fileinfo := *syscall.Stat_t{}
	if noderef {
		err := syscall.Lstat(name, fileinfo)
	} else {
		err := syscall.Stat(name, fileinfo)
	}
	sec, nsec := fileinfo.Atimespec.Unix()
	atime = time.Unix(sec, nsec)
	sec, nsec = fileinfo.Mtimespec.Unix()
	mtime = time.Unix(sec, nsec)
	return atime, mtime
}
