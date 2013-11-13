package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"os"
	"syscall"
)

var License = `License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law`

func version() error {
	fmt.Println(goopt.Suite + " " + goopt.Version)
	fmt.Println()
	fmt.Println("Copyright (C) 2013 " + goopt.Author)
	fmt.Println(License)
	os.Exit(0)
	return nil
}

func main() {
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "William Pearson"
	goopt.Version = "Unlink v0.1"
	goopt.Summary = "Uses unlink to remove FILE"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s FILE\n or:\t%s OPTION\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	switch {
	case len(os.Args) == 1:
		fmt.Println("Missing filenames")
	case len(os.Args) > 2:
		fmt.Println("Too many filenames")
	}
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	file := os.Args[1]
	if err := syscall.Unlink(file); err != nil {
		fmt.Println("Encountered an error during unlinking: %v", err)
		os.Exit(1)
	}
	return
}
