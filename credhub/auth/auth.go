// CredHub authentication strategies
package auth

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

// Provides http.Client-like interface to send authenticated requests to the server
type Auth interface {
	Do(req *http.Request) (*http.Response, error)
}

// Method is used to select an authentication strategy for CredHub.New()
//
// The server.Server provided to credhub.New() will be given to Method to construct
// the specified auth strategy.
type Method func(*server.Server) Auth
