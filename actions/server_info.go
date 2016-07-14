package actions

import (
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type ServerInfo struct {
	httpClient client.HttpClient
	config     config.Config
}
