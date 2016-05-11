package commands

import (
	"net/http"

	"fmt"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type DeleteCommand struct {
	SecretIdentifier string `short:"n" long:"name" description:"Selects the secret to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	err := config.ValidateConfig(cfg)
	if err != nil {
		return err
	}

	request := client.NewDeleteSecretRequest(cfg.ApiURL, cmd.SecretIdentifier)

	_, err = http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	fmt.Println("Secret successfully deleted")

	return nil
}
