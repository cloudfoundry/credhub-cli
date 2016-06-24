package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type GetCommand struct {
	SecretIdentifier string `short:"n" required:"yes" long:"name" description:"Selects the secret to retrieve"`
}

func (cmd GetCommand) Execute([]string) error {
	secretRepository := repositories.NewSecretRepository(client.NewHttpClient())
	config := config.ReadConfig()
	action := actions.NewSecretAction(secretRepository, config)
	secret, err := action.DoSecretAction(client.NewGetSecretRequest(config.ApiURL, cmd.SecretIdentifier), cmd.SecretIdentifier)
	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}
