// CredHub API client
package credhub

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

// CredHub server
type Config struct {
	// Url to CredHub server
	ApiUrl string

	// CA Certs in PEM format
	CaCerts []string

	// Skip certificate verification
	InsecureSkipVerify bool
}

// CredHub client to access CredHub APIs.
//
// Use New() to construct a new CredHub object, which can then interact with the CredHub api.
type CredHub struct {
	// Config provides the connection details for the CredHub server
	*Config

	// Auth provides http.Client-like Do method for authenticated requests to the CredHub server
	// Can be typecast to a specific Auth type to get additional information and functionality.
	// eg. auth.Uaa provides Logout(), Refresh(), AccessToken and RefreshToken
	auth.Auth
}

// New creates a new CredHub API client with the provided server credentials and authentication method.
// See the auth package for supported authentication methods.
func New(conf *Config, authMethod auth.Method) *CredHub {
	credhub := &CredHub{
		Config: conf,
	}

	credhub.Auth = authMethod(credhub)

	return credhub
}

// Request sends an authenticated request to the CredHub server.
//
// The pathStr should include the full path (eg. /api/v1/data) and any query parameters.
// The request body should be marshallable to JSON, but can be left nil for GET requests.
//
// Request() is used by other CredHub client methods to send authenticated requests to the CredHub server.
//
// Use Request() directly to access the CredHub server if an appropriate helper method is not available.
// For unauthenticated requests (eg. /health), use Config.Client() instead.
func (ch *CredHub) Request(method string, pathStr string, body interface{}) (*http.Response, error) {
	url, _ := url.Parse(ch.Config.ApiUrl)
	url.Path = pathStr

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url.String(), bytes.NewReader(jsonBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return ch.Auth.Do(req)
}
