package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var License = `License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law`

func version() error {
	fmt.Println(goopt.Suite + " " + goopt.Version)
	fmt.Println()
	fmt.Println("Copyright (C) 2013 " + goopt.Author)
	fmt.Println(License)
	os.Exit(0)
	return nil
}

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

func promptBeforeOverwrite(filename string) bool {
	prompt := "Overwrite " + filename + "?"
	var response string
	trueresponse := "yes"
	falseresponse := "no"
	for {
		fmt.Print(prompt)
		fmt.Scanln(&response)
		response = strings.ToLower(response)
		if strings.Contains(trueresponse, response) {
			return true
		} else if strings.Contains(falseresponse, response) || response == "" {
			return false
		}
	}
}

func main() {
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "William Pearson"
	goopt.Version = "Mv v0.1"
	goopt.Summary = "Move (rename) each SOURCE to DEST"
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
	/*recurse := goopt.Flag([]string{"-r", "-R", "--recurse"}, nil, "Recursively copy files from SOURCE to TARGET", "")*/
	goopt.OptArg([]string{"-S", "--suffix"}, "SUFFIX", "Override the usual backup suffix", setBackupSuffix)
	symlink := goopt.Flag([]string{"-s", "--symbolic-link"}, nil, "Make symlinks instead of copying", "")
	goopt.OptArg([]string{"-t", "--target"}, "TARGET", "Set the target with a flag instead of at the end", setTarget)
	update := goopt.Flag([]string{"-u", "--update"}, nil, "Move only when DEST is missing or older than SOURCE", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	if len(goopt.Args) < 2 {
		fmt.Println(goopt.Usage())
		os.Exit(1)
	}
	j := 0
	if target == "" {
		target = goopt.Args[len(goopt.Args)-1]
		j = 1
	}
	/* Recursive shit to come */
	var sources []string
	for i := range goopt.Args[j:] {
		sources = append(sources, goopt.Args[i])
	}
	var destinfo os.FileInfo
	var err error
	if *nodereference {
		destinfo, err = os.Lstat(target)
	} else {
		destinfo, err = os.Stat(target)
	}
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error trying to get info to check if DEST is a directory:", err)
		os.Exit(1)
	}
	isadir := err == nil && destinfo.IsDir()
	if (len(goopt.Args) > 2 || (target != goopt.Args[len(goopt.Args)-1] && len(goopt.Args) > 1)) && !isadir {
		fmt.Println("Too many arguments for non-directory destination")
		os.Exit(1)
	}
	for i := range sources {
		dest := target
		if isadir {
			dest = dest + string(os.PathSeparator) + filepath.Base(sources[i])
		}
		if *nodereference {
			destinfo, err = os.Lstat(target)
		} else {
			destinfo, err = os.Stat(target)
		}
		exist := !os.IsNotExist(err)
		newer := true
		if err != nil && exist {
			fmt.Println("Error trying to get info on target:", err)
			os.Exit(1)
		}
		if *update && exist {
			var srcinfo os.FileInfo
			if *nodereference {
				srcinfo, err = os.Lstat(sources[i])
			} else {
				srcinfo, err = os.Stat(sources[i])
			}
			if err != nil {
				fmt.Println("Error trying to get mod time on SRC:", err)
				os.Exit(1)
			}
			newer = srcinfo.ModTime().After(destinfo.ModTime())
		}
		if newer {
			promptres := true
			if exist {
				promptres = !*noclobber
				if *prompt {
					promptres = promptBeforeOverwrite(dest)
				}
				if promptres && *backup {
					err = os.Rename(dest, dest+backupsuffix)
					if err != nil {
						fmt.Println("Error while backing up", dest, "to", dest+backupsuffix, ":", err)
						os.Exit(1)
					}
				}
			}
			if promptres {
				switch {
				case *hardlink:
					if err := os.Link(sources[i], dest); err != nil {
						fmt.Println("Error while linking", dest, "to", sources[i], ":", err)
						defer os.Exit(1)
					} else if *verbose {
						fmt.Println("Linked", dest, "to", sources[i], ":", err)
					}
				case *symlink:
					if err := os.Symlink(sources[i], dest); err != nil {
						fmt.Println("Error while linking", dest, "to", sources[i], ":", err)
						defer os.Exit(1)
					} else if *verbose {
						fmt.Println("Linked", dest, "to", sources[i], ":", err)
					}
				default:
					source, err := os.Open(sources[i])
					if err != nil {
						fmt.Println("Error while opening source file,", sources[i], ":", err)
						defer os.Exit(1)
						continue
					}
					filebuf, err := ioutil.ReadAll(source)
					if err != nil {
						fmt.Println("Error while reading source file,", sources[i], "for copying:", err)
						defer os.Exit(1)
						continue
					}
					destfile, err := os.Create(dest)
					if err != nil {
						fmt.Println("Error while creating destination file,", dest, ":", err)
						defer os.Exit(1)
						continue
					}
					destfile.Write(filebuf)
					if *verbose {
						fmt.Println(sources[i], "copied to", dest)
					}
				}
			}
		}
	}
	return
}
