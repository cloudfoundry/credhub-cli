package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
)

type LogoutCommand struct {
}

func (cmd LogoutCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	RevokeTokenIfNecessary(cfg)
	MarkTokensAsRevokedInConfig(&cfg)
	if err := config.WriteConfig(cfg); err != nil {
		return err
	}
	fmt.Println("Logout Successful")
	return nil
}

func RevokeTokenIfNecessary(cfg config.Config) error {
	credhubClient, err := credhub.New(cfg.ApiURL, credhub.CaCerts(cfg.CaCerts...), credhub.SkipTLSValidation(cfg.InsecureSkipVerify))
	if err != nil {
		return err
	}

	uaaClient := uaa.Client{
		AuthURL: cfg.AuthURL,
		Client:  credhubClient.Client(),
	}

	uaaClient.RevokeToken(cfg.AccessToken)
	return nil
}

func MarkTokensAsRevokedInConfig(cfg *config.Config) {
	cfg.AccessToken = "revoked"
	cfg.RefreshToken = "revoked"
}
