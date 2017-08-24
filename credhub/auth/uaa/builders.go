package uaa

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

// Provides a Builder for a UAA authentication strategy using password grant type
func UaaPassword(clientId, clientSecret, username, password string) auth.Builder {
	return Uaa(clientId, clientSecret, username, password, "", "")
}

// UaaClientCredentialGrant provides a Builder for a UAA authentication strategy
// using client_credentials grant type
func UaaClientCredentials(clientId, clientSecret string) auth.Builder {
	return Uaa(clientId, clientSecret, "", "", "", "")
}

// Provides a Builder for a UAA authentication strategy using existing tokens
func Uaa(clientId, clientSecret, username, password, accessToken, refreshToken string) auth.Builder {
	return func(config auth.Config) (auth.Strategy, error) {
		httpClient := config.Client()
		authUrl, err := config.AuthURL()

		if err != nil {
			return nil, err
		}

		uaaClient := Client{
			AuthURL: authUrl,
			Client:  httpClient,
		}

		oauth := &auth.OAuthStrategy{
			Username:     username,
			Password:     password,
			ClientId:     clientId,
			ClientSecret: clientSecret,
			ApiClient:    httpClient,
			OAuthClient:  &uaaClient,
		}

		oauth.SetTokens(accessToken, refreshToken)

		return oauth, nil
	}
}
