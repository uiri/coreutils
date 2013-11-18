package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
)

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
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) != 2 {
		coreutils.PrintUsage()
	}
	file1 := goopt.Args[0]
	file2 := goopt.Args[1]
	if err := os.Link(file1, file2); err != nil {
		fmt.Fprintf(os.Stderr, "Encountered an error during linking: %v\n", err)
		os.Exit(1)
	}
	return
}
