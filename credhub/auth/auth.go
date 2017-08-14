// CredHub authentication strategies
package auth

import (
	"net/http"
)

// Provides http.Client-like interface to send authenticated requests to the server
type Auth interface {
	Do(req *http.Request) (*http.Response, error)
}

type ServerConfig interface {
	AuthUrl() (string, error)
	Client() *http.Client
}

// Builder is used to select an authentication strategy for CredHub.New()
//
// The ServerConfig provided to credhub.New() will be given to Builder to construct
// the specified auth strategy.
type Builder func(config ServerConfig) Auth
