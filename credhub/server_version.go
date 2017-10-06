package credhub

import version "github.com/hashicorp/go-version"

func (ch *CredHub) ServerVersion() (*version.Version, error) {
	if ch.cachedServerVersion == "" {
		info, err := ch.Info()
		if err != nil {
			return nil, err
		}
		ch.cachedServerVersion = info.App.Version
	}
	return version.NewVersion(ch.cachedServerVersion)
}
