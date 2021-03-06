// +build !plan9,!windows

package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"syscall"
)

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Unlink"
	goopt.Summary = "Uses unlink to remove FILE"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s FILE\n or:\t%s OPTION\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) != 1 {
		coreutils.PrintUsage()
	}
	file := goopt.Args[0]
	if err := syscall.Unlink(file); err != nil {
		fmt.Fprintf(os.Stderr, "Encountered an error during unlinking: %v\n", err)
		os.Exit(1)
	}
	return
}
