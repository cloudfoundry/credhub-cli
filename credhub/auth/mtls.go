package auth

import (
	"net/http"
)

// Mutual TLS authentication strategy
//
// When a MutualTls auth.Builder (eg. MutualTlsCertificate()) is provided to credhub.New(),
// CredHub will use this MutualTls.Do() to send authenticated requests to CredHub.
type MutualTls struct {
	Certificate string
}

// Provides http.Client-like interface to send requests authenticated with MutualTLS
func (a *MutualTls) Do(http.Request) (http.Response, error) {
	panic("Not implemented")
}
