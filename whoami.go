package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"os"
	"os/user"
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
	goopt.Version = "Whoami v0.1"
	goopt.Summary = "Prints username of current user"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage: %s OPTION\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user: %v", err)
		return
	}
	fmt.Println(currentUser.Username)
}
