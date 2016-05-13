package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/cm-cli/commands"
)

func main() {
	if os.Args[1] == "--version" {
		os.Args[1] = "version"
	}

	parser := flags.NewParser(&commands.CM, flags.HelpFlag)

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
