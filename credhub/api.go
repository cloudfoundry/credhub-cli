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
// Use New() to construct a new CredHub object, which can then interact with the CredHub api.
type CredHub struct {
	ApiURL string

	// Strategy provides http.Client-like Do method for authenticated requests to the CredHub server
	// Can be typecast to a specific Strategy type to get additional information and functionality.
	// eg. auth.OAuthStrategy provides Logout(), Refresh(), AccessToken and RefreshToken
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

// New creates a new CredHub API client with the provided server credentials and authentication method.
// See the auth package for supported authentication methods.
func New(addr string, options ...func(*CredHub) error) (*CredHub, error) {
	baseURL, err := url.Parse(addr)

	if err != nil {
		return nil, err
	}

	credhub := &CredHub{
		ApiURL:  addr,
		baseURL: baseURL,
	}

	for _, option := range options {
		if err := option(credhub); err != nil {
			return nil, err
		}
	}

	if credhub.baseURL.Scheme == "https" {
		credhub.defaultClient = httpsClient(credhub.insecureSkipVerify, credhub.caCerts)
	} else {
		credhub.defaultClient = httpClient()
	}

	if credhub.Auth != nil {
		return credhub, nil
	}

	if credhub.authBuilder == nil {
		credhub.Auth = credhub.defaultClient
		return credhub, nil
	}

	credhub.Auth, err = credhub.authBuilder(credhub)

	if err != nil {
		return nil, err
	}

	return credhub, nil
}
