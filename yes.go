package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"os"
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
	goopt.Version = "Yes v0.1"
	goopt.Summary = "Loops forever printing out a string or 'y'"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage of %s:\n\t   %s STRING\n\tor %s OPTION\n", os.Args[0], os.Args[0], os.Args[0]) +
			goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	repeat := "y"
	for i := range os.Args[1:] {
		if i > 0 {
			repeat = repeat + " " + os.Args[i+1]
		} else {
			repeat = os.Args[i+1]
		}
	}
	for {
		fmt.Println(repeat)
	}
}
