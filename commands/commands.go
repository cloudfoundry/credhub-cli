package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

type NeedsClient struct {
	client *credhub.CredHub
}

func (n *NeedsClient) SetClient(client *credhub.CredHub) {
	n.client = client
}

type NeedsConfig struct {
	config config.Config
}

func (n *NeedsConfig) SetConfig(config config.Config) {
	n.config = config
}
