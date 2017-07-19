package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"log"

	"crypto/x509"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

const TIMEOUT_SECS = 45

//go:generate counterfeiter . HttpClient

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

func NewHttpClient(cfg config.Config) *http.Client {
	parsedUrl, _ := url.Parse(cfg.ApiURL)
	if parsedUrl.Scheme == "https" {
		return newHttpsClient(cfg)
	} else {
		return newHttpClient()
	}
}

func newHttpClient() *http.Client {
	return &http.Client{Timeout: time.Second * TIMEOUT_SECS}
}

func newHttpsClient(cfg config.Config) *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify:       cfg.InsecureSkipVerify,
		PreferServerCipherSuites: true,
	}

	if !cfg.InsecureSkipVerify {
		for _, cert := range cfg.CaCerts {
			if tlsConfig.RootCAs == nil {
				tlsConfig.RootCAs = x509.NewCertPool()
			}
			ok := tlsConfig.RootCAs.AppendCertsFromPEM([]byte(cert))

			if !ok {
				log.Fatal("failed to parse root certificate")
			}
		}
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * TIMEOUT_SECS,
	}
	return client
}
