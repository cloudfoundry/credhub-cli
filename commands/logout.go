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
	cfg, _ := config.ReadConfig()
	SendLogoutIfNecessary(cfg)
	cfg = RevokedConfig(cfg)
	config.WriteConfig(cfg)
	fmt.Println("Logout Successful")
	return nil
}

func SendLogoutIfNecessary(cfg config.Config) {
	if cfg.RefreshToken != "" && cfg.RefreshToken != "revoked" {
		authRepository := repositories.NewAuthRepository(client.NewHttpClient(cfg.AuthURL), false)
		authRepository.SendRequest(client.NewTokenRevocationRequest(cfg), "logout")
	}
}

func RevokedConfig(cfg config.Config) config.Config {
	cfg.AccessToken = "revoked"
	cfg.RefreshToken = "revoked"
	return cfg
}
