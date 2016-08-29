package actions

import (
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
)

type ServerInfo struct {
	httpClient client.HttpClient
	config     config.Config
}
