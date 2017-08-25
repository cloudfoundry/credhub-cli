package credhub

import (
	"encoding/json"
	"errors"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

func (ch *CredHub) Info() (*server.Info, error) {
	response, err := ch.request("GET", "/info", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	info := &server.Info{}
	decoder := json.NewDecoder(response.Body)

	if err = decoder.Decode(&info); err != nil {
		return nil, err
	}

	return info, nil
}

// Provides the authentication server's URL
func (ch *CredHub) AuthURL() (string, error) {
	if ch.authURL != nil {
		return ch.authURL.String(), nil
	}

	info, err := ch.Info()

	if err != nil {
		return "", err
	}

	authUrl := info.AuthServer.URL

	if authUrl == "" {
		return "", errors.New("AuthURL not found")
	}

	return authUrl, nil
}
