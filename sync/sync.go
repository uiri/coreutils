// +build !windows,!plan9,!linux

package main

import "syscall"

func Sync() {
	if err := syscall.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Encountered an error during sync: %v\n", err)
		os.Exit(1)
	}
}
