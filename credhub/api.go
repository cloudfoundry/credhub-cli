// CredHub API client
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
	// Provides server information and http.Client for unauthenticated requests to the CredHub server
	Server *server.Server

	// Provides http.Client-like method for authenticated requests to the CredHub server
	// Can be typecast to a specific Auth type to get additional information and functionality.
	// eg. auth.Uaa provides Logout(), Refresh(), AccessToken and RefreshToken
	Auth *auth.Auth
}

// Sends an authenticated request to the CredHub server.
//
// The pathStr should include the full path (eg. /api/v1/data) and any query parameters.
//
// The request body should be marshallable to JSON, but can be left nil for GET requests.
//
// Request() is used by other CredHub client methods to send authenticated requests to the CredHub server.
//
// Use Request() directly to access the CredHub server if an appropriate helper method is not available.
//
// For unauthenticated requests (eg. /health), use Server.Client() instead.
func (ch CredHub) Request(method string, pathStr string, body interface{}) (http.Response, error) {
	panic("Not implemented")
}

// Creates a new CredHub API client with the provided server credentials and authentication method.
// See the auth package for supported authentication methods.
func New(server server.Server, authMethod auth.Method) CredHub {
	panic("Not implemented")
}
