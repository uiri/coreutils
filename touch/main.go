package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"strconv"
	"time"
)

var (
	atime         time.Time
	mtime         time.Time
	Nodereference *bool
)

func fromReference(rfile string) error {
	atime, mtime = GetAtimeMtime(rfile, *Nodereference)
	return nil
}

func parsePartialYear(stamp string) string {
	i, err := strconv.ParseInt(string(stamp[0]), 10, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Apparently the first character of your timestamp isn't a digit: %v\n", err)
		os.Exit(1)
	}
	if i < 6 {
		stamp = "20" + stamp
	} else if i > 6 {
		stamp = "19" + stamp
	} else if stamp[1] == '9' {
		stamp = "19" + stamp
	} else {
		stamp = "20" + stamp
	}
	return stamp
}

func fromStamp(stamp string) error {
	var layout string
	var err error
	switch len(stamp) {
	case 8:
		stamp = strconv.FormatInt(int64(mtime.Year()), 10) + stamp
		layout = "200601021504"
	case 10:
		stamp = parsePartialYear(stamp)
		layout = "200601021504"
	case 11:
		stamp = strconv.FormatInt(int64(mtime.Year()), 10) + stamp
		layout = "200601021504.05"
	case 12:
		layout = "200601021504"
	case 13:
		stamp = parsePartialYear(stamp)
		layout = "200601021504.05"
	case 15:
		layout = "200601021504.05"
	}
	tzs, _ := mtime.Zone()
	layout = layout + "MST"
	stamp = stamp + tzs
	mtime, err = time.Parse(layout, stamp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing stamp '%s': %v\n", stamp, err)
		os.Exit(1)
	}
	atime = mtime
	return nil
}

func main() {
	atime = time.Now()
	mtime = time.Now()
	goopt.Author = "William Pearson"
	goopt.Version = "Touch"
	goopt.Summary = "Change access or modification time of each FILE"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... FILE...\n", os.Args[0]) +
			goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	access := goopt.Flag([]string{"-a"}, nil, "Only change access time", "")
	modify := goopt.Flag([]string{"-m"}, nil, "Only change modification time", "")
	create := goopt.Flag([]string{"-c"}, nil, "Only change modification time", "")
	Nodereference = goopt.Flag([]string{"-h", "--no-dereference"}, []string{"--derference"}, "Affect symbolic links directly instead of dereferencing them", "Dereference symbolic links before operating on them (This is default)")
	goopt.OptArg([]string{"-r", "--reference"}, "RFILE", "Use RFILE's owner and group", fromReference)
	goopt.OptArg([]string{"-t"}, "STAMP", "Use [[CC]YY]MMDDhhmm[.ss] instead of now. Note hh is interpreted as from 00 to 23", fromStamp)
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	coreutils.Noderef = *Nodereference
	if len(goopt.Args) == 0 {
		coreutils.PrintUsage()
	}
	if !*access && !*modify {
		*access, *modify = true, true
	}
	for i := range goopt.Args {
		if !*access {
			atime, _ = GetAtimeMtime(goopt.Args[i], *Nodereference)
		}
		if !*modify {
			_, mtime = GetAtimeMtime(goopt.Args[i], *Nodereference)
		}
		if err := os.Chtimes(goopt.Args[i], atime, mtime); err != nil {
			if os.IsNotExist(err) {
				var err error
				if !*create {
					f, err := os.Create(goopt.Args[i])
					if err == nil {
						f.Close()
					}
				}
				if err == nil {
					continue
				}
			}
			fmt.Fprintf(os.Stderr, "Error touching file '%s': %v\n", goopt.Args[i], err)
		}
	}
	return
}
