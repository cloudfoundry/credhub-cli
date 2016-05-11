package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/config"
)

type ApiCommand struct {
	Server string `short:"s" long:"server" description:"API endpoint"`
}

func (cmd ApiCommand) Execute([]string) error {
	c := config.ReadConfig()

	if cmd.Server == "" {
		fmt.Println(c.ApiURL)
	} else {
		c.ApiURL = cmd.Server
		config.WriteConfig(c)
	}

	return nil
}
