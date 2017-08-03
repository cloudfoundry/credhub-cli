package server

import (
	"net/http"
)

type Server struct {
	ApiUrl             string
	InsecureSkipVerify bool
	CaCerts            []string
}

// Provides the Authentication server's URL
func (s *Server) AuthUrl() (string, error) {
	panic("Not implemented")
}

// Provides an unauthenticated http(s) client according to the Server fields
func (s *Server) Client() http.Client {
	panic("Not implemented")
}
