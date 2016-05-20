package commands

import (
	"fmt"

	"encoding/json"
	"net/http"

	"os"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/version"
)

type VersionCommand struct {
}

func (cmd VersionCommand) Execute() error {
	cfg := config.ReadConfig()

	request := client.NewInfoRequest(cfg.ApiURL)

	cmVersion := "Not Found"

	response, err := http.DefaultClient.Do(request)
	if err == nil && response.StatusCode == http.StatusOK {
		info := new(client.Info)

		decoder := json.NewDecoder(response.Body)
		decoder.Decode(info)

		cmVersion = info.App.Version
	}

	fmt.Println("CLI Version:", version.Version)
	fmt.Println("CM Version:", cmVersion)

	return nil
}

func init() {
	CM.Version = func() {
		VersionCommand{}.Execute()
		os.Exit(0)
	}
}
