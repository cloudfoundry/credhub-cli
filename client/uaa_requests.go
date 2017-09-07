package client

import (
	"bytes"
	"net/url"

	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

func NewClientCredentialsGrantTokenRequest(cfg config.Config, clientId, clientSecret string) *http.Request {
	authUrl := cfg.AuthURL + "/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Add("response_type", "token")
	data.Add("client_id", clientId)
	data.Add("client_secret", clientSecret)
	request, _ := http.NewRequest("POST", authUrl, bytes.NewBufferString(data.Encode()))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func NewRefreshTokenRequest(cfg config.Config) *http.Request {
	authUrl := cfg.AuthURL + "/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", cfg.RefreshToken)
	request, _ := http.NewRequest("POST", authUrl, bytes.NewBufferString(data.Encode()))
	request.SetBasicAuth(config.AuthClient, config.AuthPassword)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func NewAuthServerInfoRequest(cfg config.Config) (*http.Request, error) {
	requestUrl := cfg.AuthURL + "/info"
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	return request, nil
}
