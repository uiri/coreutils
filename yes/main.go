package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"strings"
)

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Yes"
	goopt.Summary = "Loops forever printing out a string or 'y'"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage of %s:\n\t   %s STRING\n\tor %s OPTION\n", os.Args[0], os.Args[0], os.Args[0]) +
			goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	repeat := strings.Join(goopt.Args, " ")
	if repeat == "" {
		repeat = "y"
	}
	for {
		fmt.Println(repeat)
	}
}
