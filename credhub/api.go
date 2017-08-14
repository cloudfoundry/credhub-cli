// CredHub API client
package credhub

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"crypto/x509"
	"errors"

	"crypto/tls"
	"time"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

// CredHub client to access CredHub APIs.
//
// Use New() to construct a new CredHub object, which can then interact with the CredHub api.
type CredHub struct {
	ApiURL string

	// Auth provides http.Client-like Do method for authenticated requests to the CredHub server
	// Can be typecast to a specific Auth type to get additional information and functionality.
	// eg. auth.Uaa provides Logout(), Refresh(), AccessToken and RefreshToken
	auth.Auth

	baseURL       *url.URL
	defaultClient *http.Client

	// CA Certs in PEM format
	caCerts *x509.CertPool

	// Skip certificate verification
	insecureSkipVerify bool

	authMethod auth.Method
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

	if credhub.authMethod != nil {
		credhub.Auth = credhub.authMethod(credhub)
	} else {
		credhub.Auth = credhub.defaultClient
	}

	return credhub, nil
}

func Auth(method auth.Method) func(*CredHub) error {
	return func(c *CredHub) error {
		c.authMethod = method
		return nil
	}
}

func CACerts(certs []string) func(*CredHub) error {
	return func(c *CredHub) error {
		c.caCerts = x509.NewCertPool()

		for _, cert := range certs {
			ok := c.caCerts.AppendCertsFromPEM([]byte(cert))
			if !ok {
				return errors.New("provided ca certs are invalid")
			}
		}

		return nil
	}
}

func SkipTLSValidation() func(*CredHub) error {
	return func(c *CredHub) error {
		c.insecureSkipVerify = true
		return nil
	}
}

// Request sends an authenticated request to the CredHub server.
//
// The pathStr should include the full path (eg. /api/v1/data) and any query parameters.
// The request body should be marshallable to JSON, but can be left nil for GET requests.
//
// Request() is used by other CredHub client methods to send authenticated requests to the CredHub server.
//
// Use Request() directly to access the CredHub server if an appropriate helper method is not available.
// For unauthenticated requests (eg. /health), use Config.Client() instead.
func (c *CredHub) Request(method string, pathStr string, body interface{}) (*http.Response, error) {
	url := *c.baseURL // clone
	url.Path = pathStr

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url.String(), bytes.NewReader(jsonBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.Auth.Do(req)
}

func (c *CredHub) request(method string, path string, body io.Reader) (*http.Response, error) {
	client := c.Client()

	url := *c.baseURL // clone
	url.Path = path

	request, _ := http.NewRequest(method, url.String(), body)

	return client.Do(request)
}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 45,
	}
}

func httpsClient(insecureSkipVerify bool, rootCAs *x509.CertPool) *http.Client {
	client := httpClient()

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:       insecureSkipVerify,
			PreferServerCipherSuites: true,
			RootCAs:                  rootCAs,
		},
	}

	return client
}
