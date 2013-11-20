package main

import "syscall"

func Sync() {
	syscall.Sync()
}
