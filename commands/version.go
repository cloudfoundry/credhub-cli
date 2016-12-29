package commands

import (
	"fmt"

	"os"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/version"
)

func PrintVersion() error {
	cfg := config.ReadConfig()

	cmVersion := "Not Found"
	cmInfo, err := actions.NewInfo(client.NewHttpClient(cfg), cfg).GetServerInfo()
	if err == nil {
		cmVersion = cmInfo.App.Version
	}

	fmt.Println("CLI Version:", version.Version)
	fmt.Println("API Version:", cmVersion)

	return nil
}

func init() {
	CM.Version = func() {
		PrintVersion()
		os.Exit(0)
	}
}
