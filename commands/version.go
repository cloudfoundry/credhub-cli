package commands

import (
	"fmt"

	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/version"
)

func PrintVersion() error {
	cfg := config.ReadConfig()

	credHubServerVersion := "Not Found"
	cmInfo, err := actions.NewInfo(client.NewHttpClient(cfg), cfg).GetServerInfo()
	if err == nil {
		credHubServerVersion = cmInfo.App.Version
	}

	fmt.Println("CLI Version:", version.Version)
	fmt.Println("Server Version:", credHubServerVersion)

	return nil
}

func init() {
	CM.Version = func() {
		PrintVersion()
		os.Exit(0)
	}
}
