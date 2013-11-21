package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"path/filepath"
)

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Ln"
	goopt.Summary = "Make a LINK (whose name is optionally specified) to TARGET or make a link in DIRECTORY to each TARGET."
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... TARGET (LINK)\n or:\t%s [OPTION]... TARGET... DIRECTORY\n or:\t%s [OPTION]... -t DIRECTORY TARGET...\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	backup := goopt.Flag([]string{"-b", "--backup"}, nil, "Backup files before overwriting", "")
	prompt := goopt.Flag([]string{"-i", "--interactive"}, nil, "Prompt before removing any existing files.", "")
	force := goopt.Flag([]string{"-f", "--force"}, nil, "Remove existing files which match link names", "")
	directories := goopt.Flag([]string{"-d", "-F", "--directory"}, nil, "Allow attempts to hard link directories", "")
	nodereference := goopt.Flag([]string{"-n", "--no-dereference"}, nil, "treat LINK as a normal file if it is a symbolic link", "")
	logical := goopt.Flag([]string{"-L", "--logical"}, []string{"-P", "--physical"}, "Dereference TARGET if it is a symbolic link", "Make hard links directly to symbolic links. (this is default)")
	relative := goopt.Flag([]string{"-r", "--relative"}, nil, "Make symbolic links relative to link location", "")
	symbolic := goopt.Flag([]string{"-s", "--symbolic"}, nil, "Make symbolic links rather than hard links", "")
	goopt.OptArg([]string{"-S", "--suffix"}, "SUFFIX", "Override the usual backup suffix", coreutils.SetBackupSuffix)
	goopt.OptArg([]string{"-t", "--target"}, "TARGET", "Set the target with a flag instead of at the end", coreutils.SetTarget)
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) < 2 {
		coreutils.PrintUsage()
	}
	coreutils.Noderef = *nodereference
	i := 0
	isadir := false
	exists := false
	if coreutils.Target == "" {
		if len(goopt.Args) > 1 {
			i = 1
			coreutils.Target = goopt.Args[len(goopt.Args)-1]
		} else {
			coreutils.Target = filepath.Base(goopt.Args[0])
		}
	} else {
		isadir = true
	}
	fileinfo, err := coreutils.Stat(coreutils.Target)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error getting file info for '%s': %v\n", coreutils.Target, err)
		os.Exit(1)
	}
	if err == nil {
		exists = true
		if fileinfo.IsDir() {
			isadir = true
		}
	}
	for j := range goopt.Args[i:] {
		coreutils.Noderef = !*logical
		fileinfo, err = coreutils.Stat(goopt.Args[j])
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error getting file info for '%s': %v\n", goopt.Args[j], err)
			defer os.Exit(1)
			continue
		}
		if err == nil && !*symbolic && !*directories && fileinfo.IsDir() {
			fmt.Fprintf(os.Stderr, "Attempt to hard link a directory")
			defer os.Exit(1)
			continue
		}
		coreutils.Noderef = *nodereference
		dest := coreutils.Target
		if isadir {
			dest = dest + filepath.Base(goopt.Args[j])
			fileinfo, err = coreutils.Stat(dest)
			if err != nil && !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error getting file info for '%s': %v\n", dest, err)
				defer os.Exit(1)
				continue
			}
			if err == nil {
				exists = true
			}
		}
		promptno := true
		if exists {
			if *backup {
				coreutils.Backup(dest)
			}
			if *force {
				os.Remove(dest)
			} else if *prompt {
				promptno = coreutils.PromptFunc(dest, false)
				if promptno {
					os.Remove(dest)
				}
			}
			if !*relative {
				dest, err = filepath.Abs(dest)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error making absolute path from '%s'", dest, err)
				}
			}
		}
		if !promptno {
			exists = false
			continue
		}
		if *symbolic {
			os.Symlink(goopt.Args[j], dest)
		} else {
			os.Link(goopt.Args[j], dest)
		}
		if *verbose {
			fmt.Printf("'%s' -> '%s'\n", dest, goopt.Args[j])
		}
		exists = false
	}
	return
}
