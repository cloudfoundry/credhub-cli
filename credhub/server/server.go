// CredHub server
package server

import (
	"net/http"
)

// CredHub server
type Server struct {
	// Url to CredHub server
	ApiUrl string
	// CA Certs in PEM format
	CaCerts []string
	// Skip certificate verification
	InsecureSkipVerify bool
}

// Provides the authentication server's URL
func (s *Server) AuthUrl() (string, error) {
	panic("Not implemented yet")
}

// Provides an unauthenticated http.Client to the CredHub server
func (s *Server) Client() http.Client {
	panic("Not implemented yet")
}
