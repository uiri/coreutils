package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

var (
	target       = ""
	backupsuffix = "~"
)

func setTarget(t string) error {
	target = t
	return nil
}

func setBackupSuffix(suffix string) error {
	backupsuffix = suffix
	return nil
}

func main() {
	syscall.Umask(0)
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "William Pearson"
	goopt.Version = "Cp v0.1"
	goopt.Summary = "Copy each SOURCE to DEST"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... SOURCE(...) DEST\n or:\t%s [OPTION]... -t DEST SOURCE\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	backup := goopt.Flag([]string{"-b", "--backup"}, nil, "Backup files before overwriting", "")
	prompt := goopt.Flag([]string{"-i", "--interactive"}, nil, "Prompt before an overwrite. Override -f and -n.", "")
	noclobber := goopt.Flag([]string{"-n", "--no-clobber"}, []string{"-f", "--force"}, "Do not overwrite", "Never prompt before an overwrite")
	hardlink := goopt.Flag([]string{"-l", "--link"}, nil, "Make hard links instead of copying", "")
	nodereference := goopt.Flag([]string{"-P", "--no-dereference"}, []string{"-L", "--dereference"}, "Never follow symlinks", "Always follow symlinks")
	/*preserve := goopt.Flag([]string{"-p", "--preserve"}, nil, "Preserve mode, ownership and timestamp attributes", "")*/
	recurse := goopt.Flag([]string{"-r", "-R", "--recurse"}, nil, "Recursively copy files from SOURCE to TARGET", "")
	goopt.OptArg([]string{"-S", "--suffix"}, "SUFFIX", "Override the usual backup suffix", setBackupSuffix)
	symlink := goopt.Flag([]string{"-s", "--symbolic-link"}, nil, "Make symlinks instead of copying", "")
	goopt.OptArg([]string{"-t", "--target"}, "TARGET", "Set the target with a flag instead of at the end", setTarget)
	update := goopt.Flag([]string{"-u", "--update"}, nil, "Move only when DEST is missing or older than SOURCE", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) < 2 {
		coreutils.PrintUsage()
	}
	coreutils.Noderef = *nodereference
	j := 0
	if target == "" {
		target = goopt.Args[len(goopt.Args)-1]
		j = 1
	}
	var sources []string
	for i := range goopt.Args[j:] {
		sources = append(sources, goopt.Args[i])
		if *recurse {
			if coreutils.Recurse(&sources) {
				defer os.Exit(1)
			}
		}
	}
	destinfo, err := coreutils.Stat(target)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error trying to get info to check if DEST is a directory: %v\n", err)
		os.Exit(1)
	}
	isadir := err == nil && destinfo.IsDir()
	if (len(goopt.Args) > 2 || (target != goopt.Args[len(goopt.Args)-1] && len(goopt.Args) > 1)) && !isadir {
		fmt.Fprintf(os.Stderr, "Too many arguments for non-directory destination")
		os.Exit(1)
	}
	for i := range sources {
		dest := target
		if sources[i] == "" {
			continue
		}
		destinfo, err := coreutils.Stat(target)
		exist := !os.IsNotExist(err)
		if err != nil && exist {
			fmt.Fprintf(os.Stderr, "Error trying to get info on target: %v\n", err)
			os.Exit(1)
		}
		srcinfo, err := coreutils.Stat(sources[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error trying to get mod time on SRC: %v\n", err)
			os.Exit(1)
		}
		mkdir := false
		if srcinfo.IsDir() {
			if isadir {
				dest = dest + string(os.PathSeparator) + sources[i]
			}
			if !*recurse {
				fmt.Printf("Skipping directory %s\n", sources[i])
				continue
			}
			mkdir = true
		} else if isadir {
			dest = dest + string(os.PathSeparator) + filepath.Base(sources[i])
		}
		newer := true
		if *update && exist {
			newer = srcinfo.ModTime().After(destinfo.ModTime())
		}
		if !newer {
			continue
		}
		promptres := true
		if exist {
			promptres = !*noclobber
			if *prompt {
				promptres = coreutils.PromptFunc(dest, false)
			}
			if promptres && *backup {
				if err = os.Rename(dest, dest+backupsuffix); err != nil {
					fmt.Fprintf(os.Stderr, "Error while backing up '%s' to '%s': %v\n", dest, dest+backupsuffix, err)
					defer os.Exit(1)
					continue
				}
			}
		}
		if !promptres {
			continue
		}
		switch {
		case mkdir:
			if err = os.Mkdir(dest, coreutils.Mode); err != nil {
				fmt.Fprintf(os.Stderr, "Error while making directory '%s': %v\n", dest, err)
				defer os.Exit(1)
				continue
			}
			if *verbose {
				fmt.Printf("Copying directory '%s' to '%s'\n", sources[i], dest)
			}
		case *hardlink:
			if err := os.Link(sources[i], dest); err != nil {
				fmt.Fprintf(os.Stderr, "Error while linking '%s' to '%s': %v\n", dest, sources[i], err)
				defer os.Exit(1)
				continue
			}
			if *verbose {
				fmt.Printf("Linked '%s' to '%s'\n", dest, sources[i])
			}
		case *symlink:
			if err := os.Symlink(sources[i], dest); err != nil {
				fmt.Fprintf(os.Stderr, "Error while linking '%s' to '%s': %v\n", dest, sources[i], err)
				defer os.Exit(1)
				continue
			}
			if *verbose {
				fmt.Printf("Symlinked '%s' to '%s'\n", dest, sources[i])
			}
		default:
			source, err := os.Open(sources[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while opening source file '%s': %v\n", sources[i], err)
				defer os.Exit(1)
				continue
			}
			filebuf, err := ioutil.ReadAll(source)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while reading source file, '%s', for copying: %v\n", sources[i], err)
				defer os.Exit(1)
				continue
			}
			destfile, err := os.Create(dest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while creating destination file, '%s': %v\n", dest, err)
				defer os.Exit(1)
				continue
			}
			destfile.Write(filebuf)
			if *verbose {
				fmt.Printf("'%s' copied to '%s'\n", sources[i], dest)
			}
		}
	}
	return
}
