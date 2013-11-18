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
		return nil
	}
	if when == "ALWAYS" {
		*prompteach = true
		*promptonce = false
		*force = false
		return nil
	}
	if when == "ONCE" {
		*prompteach = false
		*promptonce = true
		*force = false
	}
	return nil
}

func promptBeforeRemove(filename string, remove bool) bool {
	var prompt string
	if remove {
		prompt = "Remove " + filename + "?"
	} else {
		prompt = "Recurse into " + filename + "?"
	}
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
	emptydir := goopt.Flag([]string{"-d", "--dir"}, nil, "Remove empty directories", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	promptno := true
	var filenames []string
	var dirnames []string
	if len(goopt.Args) == 0 {
		fmt.Println(goopt.Usage())
		os.Exit(1)
	}
	for i := range goopt.Args {
		fileinfo, err := os.Lstat(goopt.Args[i])
		if err != nil {
			fmt.Println("Error getting file info,", err)
			defer os.Exit(1)
			continue
		}
		if fileinfo.IsDir() {
			dirnames = append(dirnames, goopt.Args[i])
		}
		filenames = append(filenames, goopt.Args[i])
	}
	i := 0
	l := len(filenames)
	rec := *recurse
	for rec {
		rec = false
		for j := range filenames[i:] {
			fileinfo, err := os.Lstat(filenames[i+j])
			if err != nil {
				fmt.Println("Error getting file info,", err)
				defer os.Exit(1)
				continue
			}
			if !fileinfo.IsDir() {
				continue
			}
			dirnames = append(dirnames, filenames[i+j])
			if *preserveroot && filenames[i+j] == "/" {
				continue
			}
			if *prompteach || *promptonce {
				promptno = promptBeforeRemove(filenames[i+j], false)
			}
			filelisting, err := ioutil.ReadDir(filenames[i+j])
			if err != nil && !*force {
				fmt.Println("Could not recurse into", filenames[i+j], ":", err)
				defer os.Exit(1)
				continue
			}
			if len(filelisting) == 0 {
				continue
			}
			rec = true
			for h := range filelisting {
				filenames = append(filenames, filenames[i+j]+string(os.PathSeparator)+filelisting[h].Name())
			}
		}
		i = l
		l = len(filenames)
	}
	l--
	for i := range filenames {
		isadir := false
		if *prompteach || *promptonce && (l-i)%3 == 1 {
			promptno = promptBeforeRemove(filenames[l-i], true)
		}
		for j := range dirnames {
			if filenames[l-i] == dirnames[j] {
				isadir = true
				break
			}
		}
		if !promptno {
			continue
		}
		if !*emptydir && !*recurse && isadir {
			fmt.Println("Could not remove", filenames[l-i], ": Is a directory")
			defer os.Exit(1)
			continue
		}
		if *verbose {
			fmt.Println("Removing", filenames[l-i])
		}
		err := os.Remove(filenames[l-i])
		if err != nil && !*force {
			fmt.Println("Could not remove", filenames[l-i], ":", err)
			defer os.Exit(1)
		}
	}
	return
}
