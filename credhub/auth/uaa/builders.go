package uaa

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"

// Provides a Builder for a UAA authentication strategy using password grant type
func PasswordGrantBuilder(clientId, clientSecret, username, password string) auth.Builder {
	return func(config auth.Config) (auth.Strategy, error) {
		httpClient := config.Client()
		authUrl, err := config.AuthUrl()

		if err != nil {
			return nil, err
		}

		uaaClient := Client{
			AuthUrl: authUrl,
			Client:  httpClient,
		}
		return &auth.OAuthStrategy{
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
func ClientCredentialsGrantBuilder(clientId, clientSecret string) auth.Builder {
	return func(config auth.Config) (auth.Strategy, error) {
		panic("Not implemented yet")
	}
}

// Provides a Builder for a UAA authentication strategy using existing tokens
//
// Only use this for a sessions created with password grant type
// For existing sessions created with a client_credentials grant type, use OAuthClient()
func SessionBuilder(clientId, clientSecret, accessToken, refreshToken string) auth.Builder {
	return func(config auth.Config) (auth.Strategy, error) {
		panic("Not implemented")
	}
}
