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

func frown(s string) {
	fmt.Println(os.Args[0] + ": " + s)
	fmt.Println("Try '" + os.Args[0] + " --help' for more information.")
	os.Exit(1)
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		n, interr := strconv.ParseFloat(s, 64)
		if interr != nil {
			frown("invalid time interval ‘" + s + "’")
		}
		d = time.Duration(n) * time.Second
	}
	return d
}

func main() {
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "Aaron Muir Hamilton"
	goopt.Version = "Sleep v0.1"
	goopt.Summary = "Pause for NUMBER seconds. SUFFIX may be 's' for seconds (the default), 'm' for minutes, or 'h' for hours. NUMBER may be either an integer or a floating point number."
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage of %s:\n\t   %s NUMBER[SUFFIX]\n\tor %s OPTION\n", os.Args[0], os.Args[0], os.Args[0]) +
			goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	if len(os.Args) == 1 {
		frown("missing operand")
	}
	var d time.Duration
	for i := range os.Args[1:] {
		d += parseDuration(os.Args[i+1])
	}
	time.Sleep(d)
}
