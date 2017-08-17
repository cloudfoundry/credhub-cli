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

func AuthBuilder(method auth.Builder) func(*CredHub) error {
	return func(c *CredHub) error {
		c.authBuilder = method
		return nil
	}
}

func Auth(strategy auth.Strategy) func(*CredHub) error {
	return func(c *CredHub) error {
		c.Auth = strategy
		return nil
	}
}
func AuthURL(authURL string) func(*CredHub) error {
	return func(c *CredHub) error {
		var err error
		c.authURL, err = url.Parse(authURL)
		return err
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
func (c *CredHub) Request(method string, pathStr string, query url.Values, body interface{}) (*http.Response, error) {
	u := *c.baseURL // clone
	u.Path = pathStr
	u.RawQuery = query.Encode()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(jsonBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Auth.Do(req)

	if err != nil {
		return resp, err
	}

	if err := c.checkForServerError(resp); err != nil {
		return nil, err
	}

	return resp, err
}

func (c *CredHub) request(method string, path string, body io.Reader) (*http.Response, error) {
	client := c.Client()

	url := *c.baseURL // clone
	url.Path = path

	request, _ := http.NewRequest(method, url.String(), body)

	resp, err := client.Do(request)

	if err != nil {
		return resp, err
	}

	if err := c.checkForServerError(resp); err != nil {
		return nil, err
	}

	return resp, err
}

func (c *CredHub) checkForServerError(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)

		respErr := &Error{}

		if err := dec.Decode(respErr); err != nil {
			return err
		}

		return respErr
	}

	return nil
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
