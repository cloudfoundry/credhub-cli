package auth

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
)

func Noop(config Config) (Strategy, error) {
	return &NoopStrategy{config.Client()}, nil
}

// Provides a Builder for MutualTLSStrategy authentication strategy
func MutualTLS(certificate string) Builder {
	return func(config Config) (Strategy, error) {
		panic("Not implemented")
	}
}

// Provides a Builder for a UAA authentication strategy using password grant type
func UaaPassword(clientId, clientSecret, username, password string) Builder {
	return Uaa(clientId, clientSecret, username, password, "", "")
}

// UaaClientCredentialGrant provides a Builder for a UAA authentication strategy
// using client_credentials grant type
func UaaClientCredentials(clientId, clientSecret string) Builder {
	return Uaa(clientId, clientSecret, "", "", "", "")
}

// Provides a Builder for a UAA authentication strategy using existing tokens
func Uaa(clientId, clientSecret, username, password, accessToken, refreshToken string) Builder {
	return func(config Config) (Strategy, error) {
		httpClient := config.Client()
		authUrl, err := config.AuthURL()

		if err != nil {
			return nil, err
		}

		uaaClient := uaa.Client{
			AuthURL: authUrl,
			Client:  httpClient,
		}

		oauth := &OAuthStrategy{
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
