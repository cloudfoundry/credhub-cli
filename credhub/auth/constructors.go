package auth

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/server"

// Provides a constructor for MutualTls authentication strategy
func MutualTlsCertificate(certificate string) Method {
	return func(s *server.Config) Auth {
		panic("Not implemented")
	}
}

// Provides a constructor for a UAA authentication strategy using password grant
func UaaPasswordGrant(clientId, clientSecret, username, password string) Method {
	return func(s *server.Config) Auth {
		panic("Not implemented")
		// Create an http(s) client out of server.ApiUrl
		// Hit apiUrl/info endpoint to get AuthUrl
		// Authenticate against AuthUrl based off of grant type
		// Populate Auth with access tokens
		// Populate UaaClient and ApiClient accordingly
	}
}

// UaaClientCredentialGrant provides a constructor for a UAA authentication strategy
// using client credential grant.
func UaaClientCredentialGrant(clientId, clientSecret string) Method {
	return func(s *server.Config) Auth {
		panic("Not implemented yet")
	}
}

// Provides a constructor for a UAA authentication strategy using existing tokens
//
// Only use this for a sessions created with password grant
// For existing sessions created with a client credential grant, use UaaClient()
func UaaSession(clientId, clientSecret, accessToken, refreshToken string) Method {
	return func(s *server.Config) Auth {
		panic("Not implemented")
	}
}
