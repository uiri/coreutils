package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"os"
	"time"
	"strconv"
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
	goopt.Author = "Aaron Muir Hamilton"
	goopt.Version = "Sleep v0.1"
	goopt.Summary = "Waits for a duration before continuing."
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage of %s:\n\t   %s STRING\n\tor %s OPTION\n", os.Args[0], os.Args[0], os.Args[0]) +
			goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	duration, err := time.ParseDuration(os.Args[1])
	if err != nil {
	   number, interr := strconv.ParseFloat(os.Args[1], 64)
	   if interr != nil {
	      fmt.Println(err)
	      os.Exit(0)
	   }
	   duration := time.Duration(number) * time.Second
	   time.Sleep(duration)
	   os.Exit(0)
	}
	time.Sleep(duration)
}
