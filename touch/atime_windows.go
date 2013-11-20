package main

import (
	"github.com/uiri/coreutils"
	"os"
	"syscall"
	"time"
)

func GetAtimeMtime(name string, noderef bool) (atime, mtime time.Time) {
	coreutils.Noderef = noderef
	fileinfo, err := coreutils.Stat(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting fileinfo for '%s': %v\n", name, err)
		os.Exit(1)
	}
	nsecs = fileinfo.Sys().LastAccessTime.Nanoseconds()
	secs = nsecs / uint64(1000000000)
	nsecs = nsecs % uint64(1000000000)
	atime = time.Unix(secs, nsecs)
	nsecs = fileinfo.Sys().LastAccessTime.Nanoseconds()
	secs = nsecs / uint64(1000000000)
	nsecs = nsecs % uint64(1000000000)
	mtime = time.Unix(secs, nsecs)
	return atime, mtime
}
