package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
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

var target = ""

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
	prompt := goopt.Flag([]string{"-i", "--interactive"}, nil, "Prompt before an overwrite. Override -f and -n.", "")
	noclobber := goopt.Flag([]string{"-n", "--no-clobber"}, []string{"-f", "--force"}, "Do not overwrite", "Never prompt before an overwrite")
	update := goopt.Flag([]string{"-u", "--update"}, nil, "Move only when DEST is missing or older than SOURCE", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	l := len(os.Args) - 1
	j := l
	for i := range os.Args[1:] {
		if os.Args[l-i][0] != '-' {
			j = l - i
			break
		}
	}
	if target == "" {
		target = os.Args[j]
	}
	destinfo, err := os.Lstat(target)
	if err != nil {
		fmt.Println("Error trying to get info to check if DEST is a directory:", err)
		os.Exit(1)
	}
	isadir := destinfo.IsDir()
	var sources []string
	for i := range os.Args[1:j] {
		if os.Args[i][0] != '-' {
			sources = append(sources, os.Args[i])
		}
		if len(sources) > 1 && !isadir {
			fmt.Println("Too many arguments for non-directory destination")
			os.Exit(1)
		}
	}
	for i := range sources {
		dest := target
		if isadir {
			dest = dest + string(os.PathSeparator) + sources[i]
		}
		destinfo, err := os.Lstat(dest)
		exist := !os.IsNotExist(err)
		newer := notexist
		if err != nil && exist {
			fmt.Println("Error trying to get :", err)
			os.Exit(1)
		}
		if *update && exist {
			srcinfo, err := os.Lstat(sources[i])
			if err != nil {
				fmt.Println("Error trying to get mod time on SRC:", err)
				os.Exit(1)
			}
			newer = srcinfo.ModTime().After(destinfo.ModTime())
		}
		if !*update {
			promptres := true
			if exist {
				/* TODO: BACKUP STUFF HERE */
				promptres = !*noclobber
				if *prompt {
					promptres = promptBeforeOverwrite(dest)
				}
			}
			if promptres {
				os.Rename(sources[i], dest)
			}
		}
	}
	return
}
