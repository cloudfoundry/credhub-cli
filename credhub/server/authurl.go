package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Provides the authentication server's URL
func (s *Server) AuthUrl() (string, error) {
	info, err := s.info()
	if err != nil {
		return "", err
	}

	authUrl := info.AuthServer.Url
	if authUrl == "" {
		return "", errors.New("AuthUrl not found")
	}

	return authUrl, nil
}

func (s *Server) request(method string, path string, body io.Reader) (*http.Response, error) {
	client, err := s.Client()
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest(method, s.ApiUrl+path, body)

	return client.Do(request)
}

func (s *Server) info() (info, error) {
	var i info
	response, err := s.request("GET", "/info", nil)
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
