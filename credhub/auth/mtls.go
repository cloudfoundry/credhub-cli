package auth

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

// Mutual TLS authentication strategy
//
// When a MutualTls auth.Method (eg. MutualTlsCertificate()) is provided to credhub.New(),
// CredHub will use this MutualTls.Do() to send authenticated requests to CredHub.
type MutualTls struct {
	Certificate string
}

// Provides http.Client-like interface to send requests authenticated with MutualTLS
func (a *MutualTls) Do(http.Request) (http.Response, error) {
	panic("Not implemented")
}

// Provides a constructor for MutualTls authentication strategy
func MutualTlsCertificate(certificate string) Method {
	return func(s *server.Server) Auth {
		panic("Not implemented")
	}
}
