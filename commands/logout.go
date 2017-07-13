package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type LogoutCommand struct {
}

func (cmd LogoutCommand) Execute([]string) error {
	api.Logout()

	cfg := config.ReadConfig()
	cfg.MarkTokensAsRevoked()
	config.WriteConfig(cfg)

	fmt.Println("Logout Successful")
	return nil
}
