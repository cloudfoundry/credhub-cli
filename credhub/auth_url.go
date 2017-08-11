package credhub

import (
	"errors"
)

// Provides the authentication server's URL
func (c *CredHub) AuthUrl() (string, error) {
	info, err := c.Info()

	if err != nil {
		return "", err
	}

	authUrl := info.AuthServer.Url

	if authUrl == "" {
		return "", errors.New("AuthUrl not found")
	}

	return authUrl, nil
}
