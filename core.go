package coreutils

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var License = `License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law`

var Mode = os.FileMode(uint32(0755))
var Noderef = true
var Preserveroot = false
var Silent = false
var Prompt = false

var PromptFunc = func(filename string, ignored bool) bool {
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
		}
		if strings.Contains(falseresponse, response) || response == "" {
			return false
		}
	}
}

func Version() error {
	fmt.Println(goopt.Suite + " " + goopt.Version)
	fmt.Println()
	fmt.Println("Copyright (C) 2013 " + goopt.Author)
	fmt.Println(License)
	os.Exit(0)
	return nil
}

func ParseMode(m string) error {
	smallend, err := strconv.ParseUint(m, 8, 32)
	if err == nil {
		Mode = os.FileMode(uint32(smallend))
		return nil
	}
	if err.Error() == fmt.Sprintf("strconv.ParseUint: parsing %q: invalid syntax", m) {
		Mode = os.FileMode(uint32(0))
		pieces := strings.Split(m, ",")
	Outer:
		for i := range pieces {
			user := uint32(0)
			group := uint32(0)
			other := uint32(0)
			subtract := false
			bitmask := uint32(0)
			for j := range pieces[i] {
				switch pieces[i][j] {
				case 'u':
					user = uint32(1)
				case 'g':
					group = uint32(1)
				case 'o':
					other = uint32(1)
				case 'a':
					other = uint32(1)
					user = uint32(1)
					group = uint32(1)
				case '-':
					subtract = true
				case 'r':
					bitmask += 4
				case 'w':
					bitmask += 2
				case 'x':
					bitmask += 1
				case '0', '1', '2', '3', '4', '5', '6', '7':
					intadd, err := strconv.ParseUint(string(pieces[i][j]), 8, 32)
					if err != nil {
						break Outer
					}
					bitmask += uint32(intadd)
				}
			}
			bitmask = uint32(0100*user*bitmask) + uint32(010*group*bitmask) + uint32(01*other*bitmask)
			if subtract {
				bitmask = ^bitmask
			}
			Mode = os.FileMode(uint32(Mode) | bitmask)
		}
		return nil
	}
	fmt.Fprintf(os.Stderr, "Error while parsing mode: %v\n", err)
	os.Exit(1)
	return err
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, goopt.Usage())
	os.Exit(1)
}

func Stat(file string) (os.FileInfo, error) {
	if Noderef {
		return os.Lstat(file)
	} else {
		return os.Stat(file)
	}
}

func Recurse(fileptr *[]string) (exit bool) {
	files := *fileptr
	exit = false
	recurse := true
	n := 1
	for recurse {
		recurse = false
		l := len(files) - n
		n = 0
		for i := range files[l:] {
			info, err := Stat(files[i+l])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting file info for '%s': %v\n", files[i+l], err)
				exit = true
				continue
			}
			if !info.IsDir() {
				continue
			}
			if Preserveroot && files[i+l] == "/" {
				continue
			}
			if Prompt && PromptFunc(files[i+l], false) {
				continue
			}
			listing, err := ioutil.ReadDir(files[i+l])
			if err != nil {
				if !Silent {
					fmt.Fprintf(os.Stderr, "Error while listing directory '%s': %v\n", files[i+l], err)
				}
				exit = true
				continue
			}
			n += len(listing)
			if len(listing) == 0 {
				continue
			}
			recurse = true
			for m := range listing {
				files = append(files, files[i+l]+string(os.PathSeparator)+listing[m].Name())
			}

		}
	}
	*fileptr = files
	return exit
}
