package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func Logout() error {
	cfg := config.ReadConfig()
	if cfg.RefreshToken != "" && cfg.RefreshToken != "revoked" {
		authRepository := repositories.NewAuthRepository(client.NewHttpClient(cfg), false)
		request, err := client.NewTokenRevocationRequest(cfg)
		if err == nil {
			authRepository.SendRequest(request, "logout")
		}
	}

	// FIXME should be handled by consumer
	cfg.MarkTokensAsRevoked()
	config.WriteConfig(cfg)

	return nil
}
