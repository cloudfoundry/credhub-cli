package main

import (
	"github.com/pivotal-cf/cm-cli/commands"
	"github.com/jessevdk/go-flags"
	"fmt"
	"os"
)

func main() {
	parser := flags.NewParser(&commands.CM, flags.HelpFlag)

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
