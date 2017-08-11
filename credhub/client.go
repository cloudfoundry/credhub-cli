package credhub

import (
	"net/http"
)

// Provides an unauthenticated http.Client to the CredHub server
func (c *CredHub) Client() *http.Client {
	return c.defaultClient
}
