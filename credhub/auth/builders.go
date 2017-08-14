package auth

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"

// Provides a Builder for MutualTLSStrategy authentication strategy
func MutualTLSCertificate(certificate string) Builder {
	return func(config Config) (Strategy, error) {
		panic("Not implemented")
	}
}

// Provides a Builder for a UAA authentication strategy using password grant type
func UaaPasswordGrant(clientId, clientSecret, username, password string) Builder {
	return func(config Config) (Strategy, error) {
		httpClient := config.Client()
		authUrl, err := config.AuthUrl()

		if err != nil {
			return nil, err
		}

		uaaClient := uaa.Client{
			AuthUrl: authUrl,
			Client:  httpClient,
		}
		return &OAuthStrategy{
			Username:     username,
			Password:     password,
			ClientId:     clientId,
			ClientSecret: clientSecret,
			ApiClient:    httpClient,
			OAuthClient:  &uaaClient,
		}, nil
	}
}

// UaaClientCredentialGrant provides a Builder for a UAA authentication strategy
// using client_credentials grant type
func UaaClientCredentialGrant(clientId, clientSecret string) Builder {
	return func(config Config) (Strategy, error) {
		panic("Not implemented yet")
	}
}

// Provides a Builder for a UAA authentication strategy using existing tokens
//
// Only use this for a sessions created with password grant type
// For existing sessions created with a client_credentials grant type, use OAuthClient()
func UaaSession(clientId, clientSecret, accessToken, refreshToken string) Builder {
	return func(config Config) (Strategy, error) {
		panic("Not implemented")
	}
}
