package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type DeleteCommand struct {
	SecretIdentifier string `short:"n" long:"name" required:"yes" description:"Selects the secret to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	config := config.ReadConfig()
	repository := repositories.NewSecretRepository(client.NewHttpClient(config.ApiURL))
	action := actions.NewAction(repository, config)

	_, err := action.DoAction(client.NewDeleteSecretRequest(config, cmd.SecretIdentifier), cmd.SecretIdentifier)

	if err == nil {
		fmt.Println("Secret successfully deleted")
	}

	return err
}
