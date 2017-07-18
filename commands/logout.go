package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type LogoutCommand struct {
}

func (cmd LogoutCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	api.NewApi(&cfg).Logout()
	config.WriteConfig(cfg)

	fmt.Println("Logout Successful")
	return nil
}
