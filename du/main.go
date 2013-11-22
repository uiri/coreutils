package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var (
	blockSize = 0
	file0from = ""
	threshold = 0
	maxdepth  = 0
	lastchar  = "\n"
)

func setBlockSize(size string) error {
	/* parse size into integer */
	return nil
}

func setFile0From(file string) error {
	file0from = file
	return nil
}

func setMaxDepth(depth string) error {
	mdepth, err := strconv.ParseInt(depth, 10, 0)
	maxdepth = int(mdepth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing depth %s: %v\n", depth, err)
		os.Exit(1)
	}
	return nil
}

func setThreshold(size string) error {
	thold, err := strconv.ParseInt(size, 10, 0)
	threshold = int(thold)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing threshold %s: %v\n", size, err)
		os.Exit(1)
	}
	return nil
}

func Stat(name string, deref, derefargs, isarg bool) syscall.Stat_t {
	finfo := new(syscall.Stat_t)
	var err error
	if deref || (derefargs && isarg) {
		err = syscall.Stat(name, finfo)
	} else {
		err = syscall.Lstat(name, finfo)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get fileinfo for '%s': %v\n", name, err)
		os.Exit(1)
	}
	return *finfo
}

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Du"
	goopt.Summary = "Estimate file sizes"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... [FILE]...\n or:\t%s [OPTION]... --files0-from=F\n", os.Args[0], os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	null := goopt.Flag([]string{"-0", "--null"}, nil, "End output with \\0 instead of \\n", "")
	all := goopt.Flag([]string{"-a", "--all"}, nil, "Write counts for all files instead of only directories", "")
	goopt.OptArg([]string{"-B", "--block-size"}, "SIZE", "Set the block size that is used when printing.", setBlockSize)
	/*bytes := goopt.Flag([]string{"-b", "--bytes"}, nil, "Equivalent to --block-size=1", "")*/
	/*total := goopt.Flag([]string{"-c", "--total"}, nil, "Add up all the sizes to create a total", "")*/
	derefargs := goopt.Flag([]string{"-D", "--dereference-args", "-H"}, nil, "Dereference symlinks if they are a commandline argument", "")
	goopt.OptArg([]string{"-d", "--max-depth"}, "N", "Print total for directory that is N or fewer levels deep", setMaxDepth)
	goopt.OptArg([]string{"--files0-from"}, "F", "Use \\0 terminated file names from file F as commandline arguments", setFile0From)
	/*human := goopt.Flag([]string{"-h", "--human-readable"}, nil, "Output using human readable suffices", "")*/
	/*kilo := goopt.Flag([]string{"-k"}, nil, "Equivalent to --block-size=1K", "")*/
	dereference := goopt.Flag([]string{"-L", "--dereference"}, []string{"-P", "--no-dereference"}, "Dereference symbolic links", "Do not dereference any symbolic links (this is default)")
	/*separate := goopt.Flag([]string{"-S", "--separate-dirs"}, nil, "Do not add subdirectories to a directory's size", "")*/
	summarize := goopt.Flag([]string{"-s", "--summarize"}, nil, "Display totals only for each argument", "")
	goopt.OptArg([]string{"-t", "--threshold"}, "SIZE", "Only include entries whose size is greater than or equal to SIZE", setThreshold)
	/* Time and Exclude options go here */
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if *null {
		lastchar = "\0000"
	}
	for i := range goopt.Args {
		fileinfo := Stat(goopt.Args[i], *dereference, *derefargs, true)
		deeper := append([]string{}, goopt.Args[i])
		coreutils.Recurse(&deeper)
		if !*summarize {
			orig := true
			startdepth := len(strings.Split(goopt.Args[i], string(os.PathSeparator)))
			l := len(deeper) - 1
			foldertosize := make(map[string]int64)
			for j := range deeper {
				foldertosize[filepath.Dir(deeper[l-j])] = 0
			}
			for j := range deeper {
				if maxdepth < (startdepth+len(strings.Split(deeper[l-j], string(os.PathSeparator)))) && maxdepth != 0 {
					continue
				}
				fileinfo = Stat(deeper[l-j], *dereference, *derefargs, orig)
				if fileinfo.Mode > 32000 {
					foldertosize[filepath.Dir(deeper[l-j])] += fileinfo.Size / 1024
					if *all {
						fmt.Printf("%d\t%s%s", fileinfo.Size, deeper[l-j], lastchar)
					}
				} else {
					fmt.Printf("%d\t%s%s", foldertosize[filepath.Base(deeper[l-j])], deeper[l-j], lastchar)
				}
			}
		} else {
			fmt.Printf("%d\t%s%s", fileinfo.Size, goopt.Args[i], lastchar)
		}
	}
	return
}
