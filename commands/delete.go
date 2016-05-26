package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	. "github.com/pivotal-cf/cm-cli/errors"
)

type DeleteCommand struct {
	SecretIdentifier string `short:"n" long:"name" required:"yes" description:"Selects the secret to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	err := config.ValidateConfig(cfg)
	if err != nil {
		return err
	}

	request := client.NewDeleteSecretRequest(cfg.ApiURL, cmd.SecretIdentifier)

	response, err := client.NewHttpClient().Do(request)

	if err != nil {
		return NewNetworkError()
	}

	if response.StatusCode == 404 {
		return NewSecretNotFoundError()
	} else if response.StatusCode == 200 {
		fmt.Println("Secret successfully deleted")
	} else {
		return NewSecretBadRequestError()
	}

	return nil
}
