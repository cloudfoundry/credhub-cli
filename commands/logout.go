package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
)

type LogoutCommand struct {
}

func (cmd LogoutCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	RevokeTokenIfNecessary(cfg)
	MarkTokensAsRevokedInConfig(&cfg)
	config.WriteConfig(cfg)
	fmt.Println("Logout Successful")
	return nil
}

func RevokeTokenIfNecessary(cfg config.Config) {
	uaaClient := uaa.Client{
		AuthURL: cfg.AuthURL,
		Client: client.NewHttpClient(cfg),
	}

	err = uaaClient.RevokeToken(cfg.AccessToken)
}

func MarkTokensAsRevokedInConfig(cfg *config.Config) {
	cfg.AccessToken = "revoked"
	cfg.RefreshToken = "revoked"
}
