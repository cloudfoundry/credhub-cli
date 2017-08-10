package auth

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"

// Provides a constructor for MutualTls authentication strategy
func MutualTlsCertificate(certificate string) Method {
	return func(config ServerConfig) Auth {
		panic("Not implemented")
	}
}

// Provides a constructor for a UAA authentication strategy using password grant
func UaaPasswordGrant(clientId, clientSecret, username, password string) Method {
	return func(config ServerConfig) Auth {
		httpClient, _ := config.Client()
		authUrl, _ := config.AuthUrl()
		uaaClient := uaa.Client{
			AuthUrl: authUrl,
			Client:  httpClient,
		}
		return &Uaa{
			Username:     username,
			Password:     password,
			ClientId:     clientId,
			ClientSecret: clientSecret,
			ApiClient:    httpClient,
			UaaClient:    &uaaClient,
		}
	}
}

// UaaClientCredentialGrant provides a constructor for a UAA authentication strategy
// using client credential grant.
func UaaClientCredentialGrant(clientId, clientSecret string) Method {
	return func(config ServerConfig) Auth {
		panic("Not implemented yet")
	}
}

// Provides a constructor for a UAA authentication strategy using existing tokens
//
// Only use this for a sessions created with password grant
// For existing sessions created with a client credential grant, use UaaClient()
func UaaSession(clientId, clientSecret, accessToken, refreshToken string) Method {
	return func(config ServerConfig) Auth {
		panic("Not implemented")
	}
}
