package credhub

import (
	"errors"
)

// Provides the authentication server's URL
func (c *CredHub) AuthURL() (string, error) {
	if c.authURL != nil {
		return c.authURL.String(), nil
	}

	info, err := c.Info()

	if err != nil {
		return "", err
	}

	authUrl := info.AuthServer.URL

	if authUrl == "" {
		return "", errors.New("AuthURL not found")
	}

	return authUrl, nil
}
