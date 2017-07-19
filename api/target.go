package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func (a *Api) Target(serverUrl string, caCerts []string, skipTlsValidation bool) (models.Info, error) {
	var credhubInfo models.Info

	a.Config.CaCerts = caCerts
	a.Config.ApiURL = serverUrl
	a.Config.InsecureSkipVerify = skipTlsValidation

	credhubInfo, err := actions.NewInfo(client.NewHttpClient(*a.Config), *a.Config).GetServerInfo()
	if err != nil {
		return credhubInfo, err
	}

	if a.Config.AuthURL != credhubInfo.AuthServer.Url {
		a.Config.AccessToken = "revoked"
		a.Config.RefreshToken = "revoked"
	}

	a.Config.AuthURL = credhubInfo.AuthServer.Url

	return credhubInfo, nil
}
