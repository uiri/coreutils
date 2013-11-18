package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	/*"io/ioutil"*/
	"os"
	/*"path/filepath"*/
	"strconv"
	"strings"
	"syscall"
)

var mode uint32

func setMode(m string) error {
	smallend, err := strconv.ParseUint(m, 8, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while parsing mode: %v\n", err)
		os.Exit(1)
	}
	mode = 1<<31 | uint32(smallend)
	return nil
}

func createParents(dir string, verbose bool) bool {
	error := false
	dirs := strings.Split(dir, "/")
	base := ""
	for i := range dirs {
		if verbose {
			fmt.Printf("Creating directory %s\n", base+dirs[i])
		}
		err := os.Mkdir(base+dirs[i], os.FileMode(mode))
		if err != nil && !os.IsExist(err) {
			fmt.Fprintf(os.Stderr, "Error while creating directory '%s': %v\n", base+dirs[i], err)
			error = true
		}
		base = base + dirs[i] + string(os.PathSeparator)
	}
	return error
}

func main() {
	mode = 1<<31 | 0755
	syscall.Umask(0)
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "William Pearson"
	goopt.Version = "Mkdir v0.1"
	goopt.Summary = "Create each DIRECTORY, if it does not already exist."
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... DIRECTORY...\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	goopt.OptArg([]string{"-m", "--mode"}, "MODE", "Set file mode permissions", setMode)
	parents := goopt.Flag([]string{"-p", "--parents"}, nil, "Make parent directories as needed, no error if existing", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each directory as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) == 0 {
		coreutils.PrintUsage()
	}
	for i := range goopt.Args {
		if *parents {
			if createParents(goopt.Args[i], *verbose) {
				defer os.Exit(1)
			}
			continue
		}
		if *verbose {
			fmt.Printf("Creating directory %s\n", goopt.Args[i])
		}
		err := os.Mkdir(goopt.Args[i], os.FileMode(mode))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating direction %s: %v\n", goopt.Args[i], err)
			defer os.Exit(1)
		}
	}
	return
}
