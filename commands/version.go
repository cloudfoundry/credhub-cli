package commands

import (
	"fmt"

	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/version"
)

func PrintVersion() error {
	cfg := config.ReadConfig()

	credHubServerVersion := "Not Found"
	credhubInfo, err := GetApiInfo(cfg.ApiURL, cfg.CaCerts, cfg.InsecureSkipVerify)
	if err == nil {
		credHubServerVersion = credhubInfo.App.Version
	}

	fmt.Println("CLI Version:", version.Version)
	fmt.Println("Server Version:", credHubServerVersion)

	return nil
}

func init() {
	CredHub.Version = func() {
		PrintVersion()
		os.Exit(0)
	}
}
