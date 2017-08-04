package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// Provides an unauthenticated http.Client to the CredHub server
func (s *Server) Client() (*http.Client, error) {
	parsedUrl, err := url.Parse(s.ApiUrl)
	if err != nil {
		return nil, err
	}

	if parsedUrl.Scheme == "https" {
		return s.httpsClient()
	} else {
		return s.httpClient()
	}
}

func (s *Server) httpClient() (*http.Client, error) {
	return &http.Client{
		Timeout: time.Second * 45,
	}, nil
}

func (s *Server) httpsClient() (*http.Client, error) {
	client, _ := s.httpClient()

	rootCAs, err := s.certPool()
	if err != nil {
		return nil, err
	}

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:       s.InsecureSkipVerify,
			PreferServerCipherSuites: true,
			RootCAs:                  rootCAs,
		},
	}

	return client, nil
}

func (s *Server) certPool() (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for _, cert := range s.CaCerts {
		ok := certPool.AppendCertsFromPEM([]byte(cert))
		if !ok {
			return nil, errors.New("Invalid certificate")
		}
	}

	return certPool, nil
}
