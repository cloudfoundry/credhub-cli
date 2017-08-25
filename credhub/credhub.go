// CredHub API client
package credhub

import (
	"net/http"
	"net/url"

	"crypto/x509"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

// CredHub client to access CredHub APIs.
//
// Use New() to construct a new CredHub object, which can then interact with the CredHub API.
type CredHub struct {
	// ApiURL is the full url to the target CredHub server
	ApiURL string

	// Auth provides an authentication Strategy for authenticated requests to the CredHub server
	// Can be type asserted to a specific Strategy type to get additional functionality and information.
	// eg. auth.OAuthStrategy provides Logout(), Refresh(), AccessToken() and RefreshToken()
	Auth auth.Strategy

	baseURL       *url.URL
	defaultClient *http.Client

	// CA Certs in PEM format
	caCerts *x509.CertPool

	// Skip certificate verification
	insecureSkipVerify bool

	authBuilder auth.Builder
	authURL     *url.URL
}
