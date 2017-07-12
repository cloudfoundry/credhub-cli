package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func Login(username string, password string, clientName string, clientSecret string, serverUrl string, caCert []string, skipTlsValidation bool) (models.Token, error) {
	var (
		token models.Token
		err   error
	)
	cfg := config.ReadConfig()

	if cfg.ApiURL == "" && serverUrl == "" {
		return token, errors.NewNoApiUrlSetError()
	}

	if len(caCert) > 0 {
		cfg.CaCert = caCert
	}

	if serverUrl != "" {
		err = ApiInfo(serverUrl, caCert, skipTlsValidation)
		if err != nil {
			return token, err
		}
		cfg = config.ReadConfig()
	}

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

	if err != nil {
		// FIXME should be handled by consumer
		Logout()
		return token, err
	}

	// FIXME should be handled by consumer
	cfg.AccessToken = token.AccessToken
	cfg.RefreshToken = token.RefreshToken
	config.WriteConfig(cfg)
	return token, err
}
