package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func (a *Api) Login(username string, password string, clientName string, clientSecret string) (models.Token, error) {
	var (
		token models.Token
		err   error
	)

	if clientName != "" || clientSecret != "" {
		token, err = actions.NewAuthToken(client.NewHttpClient(*a.Config), *a.Config).GetAuthTokenByClientCredential(clientName, clientSecret)
	} else {
		token, err = actions.NewAuthToken(client.NewHttpClient(*a.Config), *a.Config).GetAuthTokenByPasswordGrant(username, password)
	}

	if err != nil {
		return token, err
	}

	a.Config.AccessToken = token.AccessToken
	a.Config.RefreshToken = token.RefreshToken
	return token, err
}

func (a *Api) LoginWithPassword(username string, password string) (models.Token, error) {
	return a.Login(username, password, "", "")
}

func (a *Api) LoginWithClientCredentials(clientName string, clientSecret string) (models.Token, error) {
	return a.Login("", "", clientName, clientSecret)
}
