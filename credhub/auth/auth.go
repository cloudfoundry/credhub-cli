package auth

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

// Provides an authenticated http.Client utilizing the
// authentication method as supplied by the AuthOption
type Auth interface {
	Do(req *http.Request) (*http.Response, error)
}

// Provides a helper method that produces an authenticated client using
// the authentication method of choice (*)
type AuthOption func(*server.Server) Auth
