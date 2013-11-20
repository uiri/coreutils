package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"io/ioutil"
	"os"
	"strings"
)

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
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	promptno := true
	var filenames []string
	if len(goopt.Args) == 0 {
		coreutils.PrintUsage()
	}
	coreutils.Prompt = *prompteach || *promptonce
	coreutils.PromptFunc = func(filename string, remove bool) bool {
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
	for i := range goopt.Args {
		_, err := os.Lstat(goopt.Args[i])
		if *force && os.IsNotExist(err) {
			continue
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting file info for '%s': %v\n", goopt.Args[i], err)
			defer os.Exit(1)
			continue
		}
		filenames = append(filenames, goopt.Args[i])
	}
	i := 0
	rec := *recurse
	for rec {
		rec = false
		l := 0
		for j := range filenames[i:] {
			fileinfo, err := os.Lstat(filenames[i+j])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting file info for '%s': %v\n", filenames[i+j], err)
				defer os.Exit(1)
				continue
			}
			if !fileinfo.IsDir() {
				continue
			}
			if *preserveroot && filenames[i+j] == "/" {
				continue
			}
			if coreutils.Prompt {
				promptno = coreutils.PromptFunc(filenames[i+j], false)
			}
			filelisting, err := ioutil.ReadDir(filenames[i+j])
			if err != nil && !*force {
				fmt.Fprintf(os.Stderr, "Could not recurse into '%s': %v\n", filenames[i+j], err)
				defer os.Exit(1)
				continue
			}
			if len(filelisting) == 0 {
				continue
			}
			rec = true
			for h := range filelisting {
				filenames = append(filenames, filenames[i+j]+string(os.PathSeparator)+filelisting[h].Name())
				l++
			}
		}
		i += l
	}
	l := len(filenames) - 1
	for i := range filenames {
		fileinfo, err := os.Lstat(filenames[l-i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting file info for '%s': %v\n", filenames[l-i], err)
			defer os.Exit(1)
			continue
		}
		isadir := fileinfo.IsDir()
		if *prompteach || *promptonce && (l-i)%3 == 1 {
			promptno = coreutils.PromptFunc(filenames[l-i], true)
		}
		if !promptno {
			continue
		}
		if !*emptydir && !*recurse && isadir {
			fmt.Fprintf(os.Stderr, "Could not remove '%s': Is a directory\n", filenames[l-i])
			defer os.Exit(1)
			continue
		}
		if *verbose {
			fmt.Printf("Removing '%s'\n", filenames[l-i])
		}
		err = os.Remove(filenames[l-i])
		if err != nil && !(*force && os.IsNotExist(err)) {
			fmt.Fprintf(os.Stderr, "Could not remove '%s': %v\n", filenames[l-i], err)
			defer os.Exit(1)
		}
	}
	return
}
