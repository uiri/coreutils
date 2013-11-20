package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"strings"
	"syscall"
)

func createParents(dir string, verbose bool) bool {
	error := false
	dirs := strings.Split(dir, "/")
	base := ""
	for i := range dirs {
		if dirs[i] == "" {
			base = "/"
			continue
		}
		base = base + dirs[i] + string(os.PathSeparator)
		_, err := os.Stat(base)
		if err == nil || os.IsExist(err) {
			continue
		}
		if verbose {
			fmt.Printf("Creating directory %s\n", base)
		}
		err = os.Mkdir(base, coreutils.Mode)
		if err != nil && !os.IsExist(err) {
			fmt.Fprintf(os.Stderr, "Error while creating directory '%s': %v\n", base, err)
			error = true
		}
	}
	return error
}

func main() {
	coreutils.Mode = 1<<31 | coreutils.Mode
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
	goopt.OptArg([]string{"-m", "--mode"}, "MODE", "Set file mode permissions", coreutils.ParseMode)
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
		err := os.Mkdir(goopt.Args[i], coreutils.Mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating direction %s: %v\n", goopt.Args[i], err)
			defer os.Exit(1)
		}
	}
	return
}
