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
	MarkTokensAsRevokedInConfig(&cfg)
	config.WriteConfig(cfg)

	fmt.Println("Logout Successful")
	return nil
}

func MarkTokensAsRevokedInConfig(cfg *config.Config) {
	cfg.AccessToken = "revoked"
	cfg.RefreshToken = "revoked"
}
