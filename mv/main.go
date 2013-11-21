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
	goopt.Version = "Mv"
	goopt.Summary = "Move (rename) each SOURCE to DEST"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... SOURCE(...) DEST\n or:\t%s [OPTION]... -t DEST SOURCE\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	prompt := goopt.Flag([]string{"-i", "--interactive"}, nil, "Prompt before an overwrite. Override -f and -n.", "")
	noclobber := goopt.Flag([]string{"-n", "--no-clobber"}, []string{"-f", "--force"}, "Do not overwrite", "Never prompt before an overwrite")
	backup := goopt.Flag([]string{"-b", "--backup"}, nil, "Backup files before overwriting", "")
	goopt.OptArg([]string{"-S", "--suffix"}, "SUFFIX", "Override the usual backup suffix", coreutils.SetBackupSuffix)
	goopt.OptArg([]string{"-t", "--target"}, "TARGET", "Set the target with a flag instead of at the end", coreutils.SetTarget)
	update := goopt.Flag([]string{"-u", "--update"}, nil, "Move only when DEST is missing or older than SOURCE", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(goopt.Args) < 2 {
		coreutils.PrintUsage()
	}
	if coreutils.Target == "" {
		coreutils.Target = goopt.Args[len(goopt.Args)-1]
	}
	destinfo, err := os.Lstat(coreutils.Target)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error trying to get info to check if DEST is a directory: %v\n", err)
		os.Exit(1)
	}
	isadir := err == nil && destinfo.IsDir()
	if (len(goopt.Args) > 2 || (coreutils.Target != goopt.Args[len(goopt.Args)-1] && len(goopt.Args) > 1)) && !isadir {
		fmt.Fprintf(os.Stderr, "Too many arguments for non-directory destination")
		os.Exit(1)
	}
	for i := range goopt.Args[1:] {
		dest := coreutils.Target
		if isadir {
			dest = dest + string(os.PathSeparator) + filepath.Base(goopt.Args[i])
		}
		destinfo, err := os.Lstat(dest)
		exist := !os.IsNotExist(err)
		newer := true
		if err != nil && exist {
			fmt.Fprintf(os.Stderr, "Error trying to get info on target: %v\n", err)
			os.Exit(1)
		}
		if *update && exist {
			srcinfo, err := os.Lstat(goopt.Args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error trying to get mod time on SRC: %v\n", err)
				os.Exit(1)
			}
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
				coreutils.Backup(dest)
			}
		}
		if promptres {
			err = os.Rename(goopt.Args[i], dest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while moving '%s' to '%s': %v\n", goopt.Args[i], dest, err)
				defer os.Exit(1)
				continue
			}
			if *verbose {
				fmt.Printf("%s -> %s\n", goopt.Args[i], dest)
			}

		}
	}
	return
}
