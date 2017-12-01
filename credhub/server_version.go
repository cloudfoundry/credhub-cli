package credhub

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
	version "github.com/hashicorp/go-version"
)

func (ch *CredHub) ServerVersion() (*version.Version, error) {
	if ch.cachedServerVersion == "" {
		info, err := ch.Info()
		if err != nil {
			return nil, err
		}
		ch.cachedServerVersion = info.App.Version
		if ch.cachedServerVersion == "" {
			version, err := ch.getVersion()
			if err != nil {
				return nil, err
			}
			ch.cachedServerVersion = version
		}
	}
	return version.NewVersion(ch.cachedServerVersion)
}

func (ch *CredHub) getVersion() (string, error) {
	response, err := ch.request(ch.Client(), "GET", "/version", nil, nil)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	versionData := &server.VersionData{}
	decoder := json.NewDecoder(response.Body)

	if err = decoder.Decode(&versionData); err != nil {
		return "", err
	}

	return versionData.Version, nil
}
