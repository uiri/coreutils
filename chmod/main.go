package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
)

var (
	usingreference bool
)

func fromReference(rfile string) error {
	usingreference = true
	var fileinfo os.FileInfo
	var err error
	fileinfo, err = coreutils.Stat(rfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading info for '%s': %v\n", rfile, err)
		os.Exit(1)
	}
	coreutils.Mode = fileinfo.Mode()
	return nil
}

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Chmod"
	goopt.Summary = "Change file mode of each FILE to MODE\nWith reference, change file mode of each FILE to that of RFILE"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... [MODE] FILE...\n or:\t%s [OPTION]... --reference=RFILE FILE...\n", os.Args[0], os.Args[0]) +
			goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	silent := goopt.Flag([]string{"-f", "--silent", "--quiet"}, nil, "Suppress most error messages", "")
	/*changes := goopt.Flag([]string{"-c", "--changes"}, nil, "Like verbose but only report changes", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")*/
	nodereference := goopt.Flag([]string{"-h", "--no-dereference"}, []string{"--derference"}, "Affect symbolic links directly instead of dereferencing them", "Dereference symbolic links before operating on them (This is default)")
	preserveroot := goopt.Flag([]string{"--preserve-root"}, []string{"--no-preserve-root"}, "Don't recurse on '/'", "Treat '/' normally (This is default)")
	goopt.OptArg([]string{"--reference"}, "RFILE", "Use RFILE's owner and group", fromReference)
	recurse := goopt.Flag([]string{"-R", "--recursive"}, nil, "Operate recursively on files and directories", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) == 0 {
		coreutils.PrintUsage()
	}
	coreutils.Noderef = *nodereference
	coreutils.Preserveroot = *preserveroot
	if !usingreference {
		coreutils.ParseMode(goopt.Args[0])
	}
	var filestomod []string
	for i := range goopt.Args[1:] {
		filestomod = append(filestomod, goopt.Args[i+1])
		if *recurse && (!*preserveroot || goopt.Args[i+1] != "/") {
			if coreutils.Recurse(&filestomod) {
				defer os.Exit(1)
			}
		}
	}
	for i := range filestomod {
		err := os.Chmod(filestomod[i], coreutils.Mode)
		if err != nil && !*silent {
			fmt.Fprintf(os.Stderr, "Error changing mode for file '%s': %v\n", filestomod[i], err)
			defer os.Exit(1)
			continue
		}
	}
	return
}
