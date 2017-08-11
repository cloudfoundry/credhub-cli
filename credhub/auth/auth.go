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

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Method is used to select an authentication strategy for CredHub.New()
//
// The server.Config provided to credhub.New() will be given to Method to construct
// the specified auth strategy.
type Method func(config ServerConfig) Auth
