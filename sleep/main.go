package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"strconv"
	"strings"
	"time"
)

func frown(s string) {
	fmt.Fprintln(os.Stderr, os.Args[0]+": "+s)
	fmt.Fprintln(os.Stderr, "Try '"+os.Args[0]+" --help' for more information.")
	os.Exit(1)
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		if strings.HasSuffix(s, "d") {
			t := strings.Split(s, "d")[0]
			n, nerr := strconv.ParseFloat(t, 64)
			if nerr != nil {
				frown("invalid time interval ‘" + s + "’")
			}
			d = time.Duration(86400000000000 * n) // Ugly hack? seems that time.Duration is a nanosecond, and multiplying by anything just means that precision is your multiplicand.
		} else {
			n, nerr := strconv.ParseFloat(s, 64)
			if nerr != nil {
				frown("invalid time interval ‘" + s + "’")
			}
			d = time.Duration(n) * time.Second
		}

	}
	return d
}

func main() {
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "Aaron Muir Hamilton"
	goopt.Version = "Sleep v0.1"
	goopt.Summary = "Pause for NUMBER seconds. SUFFIX may be 's' for seconds (the default), 'm' for minutes, 'h' for hours, or 'd' for days. NUMBER may be either an integer or a floating point number."
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage of %s:\n\t   %s NUMBER[SUFFIX]\n\tor %s OPTION\n", os.Args[0], os.Args[0], os.Args[0]) +
			goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string { return goopt.Summary + "\n\nUnless an option is passed to it." }
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	if len(os.Args) == 1 {
		frown("missing operand")
	}
	var d time.Duration
	for i := range os.Args[1:] {
		d += parseDuration(os.Args[i+1])
	}
	time.Sleep(d)
	os.Exit(0)
}
