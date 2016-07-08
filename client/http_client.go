package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/pivotal-cf/cm-cli/config"
)

const TIMEOUT_SECS = 30

//go:generate counterfeiter . HttpClient

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

func NewHttpClient(config config.Config) *http.Client {
	parsedUrl, _ := url.Parse(config.ApiURL)
	if parsedUrl.Scheme == "https" {
		return newHttpsClient()
	} else {
		return newHttpClient()
	}
}

func newHttpClient() *http.Client {
	return &http.Client{Timeout: time.Second * TIMEOUT_SECS}
}

func newHttpsClient() *http.Client {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * TIMEOUT_SECS,
	}
	return client
}
