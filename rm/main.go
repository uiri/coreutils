package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"io/ioutil"
	"os"
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
	force      *bool
	prompteach *bool
	promptonce *bool
)

func setPrompt(when string) error {
	when = strings.ToUpper(when)
	if when == "NEVER" {
		*prompteach = false
		*promptonce = false
		*force = true
	} else if when == "ALWAYS" {
		*prompteach = true
		*promptonce = false
		*force = false
	} else if when == "ONCE" {
		*prompteach = false
		*promptonce = true
		*force = false
	}
	return nil
}

func promptBeforeRemove(filename string) bool {
	/* TODO: Fill out prompt method */
	return false
}

func main() {
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "William Pearson"
	goopt.Version = "Rm v0.1"
	goopt.Summary = "Remove each FILE"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... FILE...\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	force = goopt.Flag([]string{"-f", "--force"}, nil, "Ignore nonexistent files, don't prompt user", "")
	prompteach = goopt.Flag([]string{"-i"}, nil, "Prompt before each removal", "")
	promptonce = goopt.Flag([]string{"-I"}, nil, "Prompt before removing multiple files at once", "")
	goopt.OptArg([]string{"--interactive"}, "WHEN", "Prompt according to WHEN", setPrompt)
	/*onefs := goopt.Flag([]string{"--one-file-system"}, nil, "When -r is specified, skip directories on different filesystems", "")*/
	preserveroot := goopt.Flag([]string{"--no-preserve-root"}, []string{"--preserve-root"}, "Do not treat '/' specially", "Do not remove '/' (This is default)")
	recurse := goopt.Flag([]string{"-r", "-R", "--recursive"}, nil, "Recursively remove directories and their contents", "")
	/*emptydir := goopt.Flag([]string{"-d", "--dir"}, nil, "Remove empty directories", "")*/
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	doubledash := false
	promptno := false
	for i := range os.Args[1:] {
		if !doubledash && os.Args[i+1][0] == '-' {
			if os.Args[i+1] == "--" {
				doubledash = true
			}
			continue
		}
		filenames := []string{os.Args[i+1]}
		for j := 0; j < len(filenames); j++ {
			if *prompteach || *promptonce {
				promptno = promptBeforeRemove(filenames[j])
			}
			if *verbose {
				fmt.Printf("Removing %s\n", filenames[j])
			}
			if *recurse && (!*preserveroot || filenames[j] != "/") {
				filelisting, err := ioutil.ReadDir(filenames[j])
				if err != nil && !*force {
					fmt.Println("Could not recurse into", filenames[j], ":", err)
				} else {
					for h := range filelisting {
						filenames = append(filenames, filelisting[h].Name())
					}
				}
			}
			if !promptno {
				err := os.Remove(filenames[j])
				if err != nil && !*force {
					fmt.Println("Could not remove", filenames[j], ":", err)
				}
			}
		}
	}
	return
}
