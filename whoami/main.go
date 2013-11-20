package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/uiri/coreutils"
	"os"
	"os/user"
)

func main() {
	goopt.Author = "William Pearson"
	goopt.Version = "Whoami"
	goopt.Summary = "Prints username of current user"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage: %s OPTION\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", coreutils.Version)
	goopt.Parse(nil)
	currentUser, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current user: %v\n", err)
		return
	}
	fmt.Println(currentUser.Username)
	return
}
