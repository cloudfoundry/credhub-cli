package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type LogoutCommand struct {
}

func (cmd LogoutCommand) Execute([]string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		cfg = config.Config{}
	} else {
		if cfg.RefreshToken != "" && cfg.RefreshToken != "revoked" {
			authRepository := repositories.NewAuthRepository(client.NewHttpClient(cfg.AuthURL), false)
			authRepository.SendRequest(client.NewTokenRevocationRequest(cfg), "logout")
		}
	}

	cfg.AccessToken = "revoked"
	cfg.RefreshToken = "revoked"
	config.WriteConfig(cfg)
	fmt.Println("Logout Successful")
	return nil
}
