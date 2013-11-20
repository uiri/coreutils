package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"io/ioutil"
	"os"
	"path/filepath"
)

func removeEmptyParents(dir string, verbose, ignorefail bool) bool {
	error := false
	for {
		dir = filepath.Dir(dir)
		if len(dir) == 1 {
			break
		}
		if verbose {
			fmt.Printf("Removing directory %s\n", dir)
		}
		filelisting, err := ioutil.ReadDir(dir)
		if len(filelisting) == 0 {
			err = os.Remove(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to remove %s: %v\n", dir, err)
				error = true
			}
			continue
		}
		if !ignorefail {
			fmt.Printf("Failed to remove '%s': directory not empty\n", dir)
		}
		return true
	}
	return error
}

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Rmdir"
	goopt.Summary = "Remove each DIRECTORY if it is empty"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... DIRECTORY...\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	ignorefail := goopt.Flag([]string{"--ignore-fail-on-non-empty"}, nil,
		"Ignore each failure that is from a directory not being empty", "")
	parents := goopt.Flag([]string{"-p", "--parents"}, nil, "Remove DIRECTORY and ancestors if ancestors become empty", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each directory as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) == 0 {
		coreutils.PrintUsage()
	}
	for i := range goopt.Args {
		filelisting, err := ioutil.ReadDir(goopt.Args[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove %s: %v\n", os.Args[i+1], err)
			defer os.Exit(1)
			continue
		}
		if !*ignorefail && len(filelisting) > 0 {
			fmt.Fprintf(os.Stderr, "Failed to remove '%s' directory is non-empty\n", goopt.Args[i])
			defer os.Exit(1)
			continue
		}
		if *verbose {
			fmt.Printf("Removing directory %s\n", goopt.Args[i])
		}
		err = os.Remove(goopt.Args[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove %s: %v\n", goopt.Args[i], err)
			defer os.Exit(1)
			continue
		}
		if !*parents {
			continue
		}
		dir := goopt.Args[i]
		if dir[len(dir)-1] == '/' {
			dir = filepath.Dir(dir)
		}
		if removeEmptyParents(dir, *verbose, *ignorefail) {
			defer os.Exit(1)
		}
	}
	return
}
