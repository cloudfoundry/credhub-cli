package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func Login(username string, password string, clientName string, clientSecret string) (models.Token, error) {
	var (
		token models.Token
		err   error
	)
	cfg := config.ReadConfig()

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
		token, err = actions.NewAuthToken(client.NewHttpClient(cfg), cfg).GetAuthTokenByClientCredential(clientName, clientSecret)
	} else {
		token, err = actions.NewAuthToken(client.NewHttpClient(cfg), cfg).GetAuthTokenByPasswordGrant(username, password)
	}

	return token, err
}
