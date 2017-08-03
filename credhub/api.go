// CredHub provides methods to interact with the CredHub server.
package credhub

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

// CredHub client to access CredHub APIs.
//
// Use New() to construct a new CredHub object, which can then interact with the CredHub api.
type CredHub struct {
	// Provides server information and http.Client for unauthenticated requests to the Credhub server
	Server *server.Server

	// Provides http.Client for authenticated requests to the CredHub server
	// Can be typecast to a specific Auth type to get additional information and functionality.
	// eg. auth.Uaa provides Logout(), Refresh(), AccessToken and RefreshToken
	Auth *auth.Auth
}

// Send a request to the CredHub server. The pathStr should include the full path (eg. /api/v1/data)
// and any query parameters. The request body should be marshallable to JSON, but can be left nil
// for GET requests.
//
// Request() is used by other CredHub client methods to send authenticated requests to the CredHub server.
// When using auth.Uaa, if the access token has expired, the access token will be refreshed and the
// request will be retried.
//
// Use Request() directly to access the CredHub server if an appropriate helper method is not available.
func (ch CredHub) Request(method string, pathStr string, body interface{}) (http.Response, error) {
	panic("Not implemented")
}

// Creates a new CredHub API client with the provided server credentials and authentication method.
// See the auth package for supported authentication methods.
func New(server server.Server, authOption auth.AuthOption) CredHub {
	panic("Not implemented")
}
