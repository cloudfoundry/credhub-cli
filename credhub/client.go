package credhub

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

// Provides an unauthenticated http.Client to the CredHub server
func (c *CredHub) Client() *http.Client {
	return c.defaultClient
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
