package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"os"
	"strconv"
	"strings"
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
	goopt.Version = "Echo v0.1"
	goopt.Summary = "Echo ARGs to stdout."
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [SHORT-OPTION] ARGS...\n or:\t%s LONG-OPTION\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nValid backslash escape sequences go here."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", version)
	newline := goopt.Flag([]string{"-n"}, nil, "Don't print out a newline after ARGS", "")
	backslashescape := goopt.Flag([]string{"-e"}, []string{"-E"}, "Enable interpretation of backslash escapes", "Disable interpretation of backslash escapes")

	goopt.Parse(nil)
	argstring := strings.Join(goopt.Args, " ")
	if *backslashescape {
		argstring = fmt.Sprintf("%q", argstring)
		argstring = strings.Replace(argstring, "\\\\", "\\", -1)
		validEscapeSeqs := "\\abcefnrtv0x"
		for i := 0; i < len(argstring); i++ {
			if argstring[i] != '\\' {
				continue
			}
			if !strings.Contains(validEscapeSeqs, string(argstring[i+1])) {
				argstring = argstring[:i] + "\\" + argstring[i:]
			}
			i++
		}
		backslashstring, err := strconv.Unquote(argstring)
		if err != nil {
			fmt.Printf("Error encountered when interpreting escape sequences: %v\n", err)
			os.Exit(1)
		}
		argstring = backslashstring
	}
	if *newline {
		fmt.Print(argstring)
	} else {
		fmt.Println(argstring)
	}
	return
}
