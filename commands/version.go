package commands

import (
	"fmt"

	"os"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/version"
)

func PrintVersion() error {
	cfg, _ := config.ReadConfig()

	cmVersion := "Not Found"
	cmInfo, err := actions.NewInfo(client.NewHttpClient(cfg.ApiURL), cfg).GetServerInfo()
	if err == nil {
		cmVersion = cmInfo.App.Version
	}

	fmt.Println("CLI Version:", version.Version+" build "+version.BuildNumber)
	fmt.Println("CM Version:", cmVersion)

	return nil
}

func init() {
	CM.Version = func() {
		PrintVersion()
		os.Exit(0)
	}
}
