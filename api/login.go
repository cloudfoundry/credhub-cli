package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func (a *Api) Login(username string, password string, clientName string, clientSecret string) (models.Token, error) {
	var (
		token models.Token
		err   error
	)

	if clientName != "" || clientSecret != "" {
		if username != "" || password != "" {
			return token, errors.NewMixedAuthorizationParametersError()
		}

		if clientName == "" || clientSecret == "" {
			return token, errors.NewClientAuthorizationParametersError()
		}
	}

	if username == "" && password != "" {
		return token, errors.NewPasswordAuthorizationParametersError()
	}

	if err != nil {
		return token, err
	}

	if clientName != "" || clientSecret != "" {
		token, err = actions.NewAuthToken(client.NewHttpClient(*a.Config), *a.Config).GetAuthTokenByClientCredential(clientName, clientSecret)
	} else {
		token, err = actions.NewAuthToken(client.NewHttpClient(*a.Config), *a.Config).GetAuthTokenByPasswordGrant(username, password)
	}

	return token, err
}
