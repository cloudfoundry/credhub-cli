package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func (a *Api) Logout() error {
	if a.Config.RefreshToken != "" && a.Config.RefreshToken != "revoked" {
		authRepository := repositories.NewAuthRepository(client.NewHttpClient(*a.Config), false)
		request, err := client.NewTokenRevocationRequest(*a.Config)
		if err == nil {
			authRepository.SendRequest(request, "logout")
		}
	}

	a.Config.AccessToken = "revoked"
	a.Config.RefreshToken = "revoked"

	return nil
}
