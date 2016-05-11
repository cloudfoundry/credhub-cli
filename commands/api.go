package commands

import (
	"fmt"
	"strings"

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
		if strings.HasPrefix(cmd.Server, "http://") || strings.HasPrefix(cmd.Server, "https://") {
			c.ApiURL = cmd.Server
		} else {
			c.ApiURL = "http://" + cmd.Server
		}

		config.WriteConfig(c)

	}

	return nil
}
