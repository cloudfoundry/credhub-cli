package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

const TIMEOUT_SECS = 30

//go:generate counterfeiter . HttpClient

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

func NewHttpClient(serverUrl string) *http.Client {
	parsedUrl, _ := url.Parse(serverUrl)
	if parsedUrl.Scheme == "https" {
		return newHttpsClient()
	} else {
		return newHttpClient()
	}
}

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * TIMEOUT_SECS,
		CheckRedirect: redirectPolicyFunc,
	}
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error{
	// todo test this in!!
	req.SetBasicAuth("credhub", "") // todo move elsehwerew
	return nil
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
