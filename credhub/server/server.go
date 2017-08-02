package server

import (
	"net/http"
)

type Server struct {
	ApiUrl             string
	InsecureSkipVerify bool
	CaCerts            []string
}

func (s *Server) AuthUrl() (interface{}, error) {
	panic("Not implemented")
}

func (s *Server) Client() http.Client {
	panic("Not implemented")
}
