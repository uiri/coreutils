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
	goopt.Version = "Link v0.1"
	goopt.Summary = "Creates a link to FILE1 called FILE2"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s FILE1 FILE2\n or:\t%s OPTION\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	switch {
	case len(os.Args) == 1:
		fmt.Println("Missing filenames")
	case len(os.Args) == 2:
		fmt.Println("Missing filename after '%s'", os.Args[1])
	case len(os.Args) > 3:
		fmt.Println("Too many filenames")
	}
	if len(os.Args) != 3 {
		os.Exit(1)
	}
	file1 := os.Args[1]
	file2 := os.Args[2]
	if err := os.Link(file1, file2); err != nil {
		fmt.Println("Encountered an error during linking: %v", err)
		os.Exit(1)
	}
	return
}
