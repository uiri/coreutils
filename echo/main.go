package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"strconv"
	"strings"
)

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Echo"
	goopt.Summary = "Echo ARGs to stdout."
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [SHORT-OPTION] ARGS...\n or:\t%s LONG-OPTION\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nValid backslash escape sequences go here."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", coreutils.Version)
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
			fmt.Fprintf(os.Stderr, "Error encountered when interpreting escape sequences: %v\n", err)
			os.Exit(1)
		}
		argstring = backslashstring
	}
	nl := "\n"
	if *newline {
		nl = ""
	}
	fmt.Printf("%s%s", argstring, nl)
	return
}
