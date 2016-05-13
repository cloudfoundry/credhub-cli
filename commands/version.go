package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/version"
)

type VersionCommand struct {
}

func (cmd VersionCommand) Execute([]string) error {
	fmt.Printf("CLI Version: %s\n", version.Version)
	return nil
}
