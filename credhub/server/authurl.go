package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Provides the authentication server's URL
func (c *Config) AuthUrl() (string, error) {
	info, err := c.info()
	if err != nil {
		return "", err
	}

	authUrl := info.AuthServer.Url
	if authUrl == "" {
		return "", errors.New("AuthUrl not found")
	}

	return authUrl, nil
}

func (c *Config) request(method string, path string, body io.Reader) (*http.Response, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest(method, c.ApiUrl+path, body)

	return client.Do(request)
}

func (c *Config) info() (info, error) {
	var i info
	response, err := c.request("GET", "/info", nil)
	if err != nil {
		return i, err
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&i)
	if err != nil {
		return i, err
	}

	return i, nil
}

type info struct {
	AuthServer struct {
		Url string
	} `json:"auth-server"`
}
