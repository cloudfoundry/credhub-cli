package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/credhub-cli/commands"
)

func main() {
	debug.SetTraceback("all")
	parser := flags.NewParser(&commands.CM, flags.HelpFlag)

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
