// CredHub authentication strategies
package auth

import (
	"net/http"
)

// Provides http.Client-like interface to send authenticated requests to the server
// Modifies the request to include authentication based on the authentication strategy
type Strategy interface {
	Do(req *http.Request) (*http.Response, error)
}

// Provides client details to Builders
type Config interface {
	AuthUrl() (string, error)
	Client() *http.Client
}

// Builder constructs the auth type given a configuration
type Builder func(config Config) (Strategy, error)
