package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
)

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Sync"
	goopt.Summary = "Flush filesystem buffers"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s OPTION\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) != 0 {
		coreutils.PrintUsage()
	}
	Sync()
	return
}
