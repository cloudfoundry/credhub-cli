package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type GetCommand struct {
	SecretIdentifier string `short:"n" required:"yes" long:"name" description:"Selects the secret to retrieve"`
}

func (cmd GetCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	repository := repositories.NewSecretRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, cfg)
	secret, err := action.DoAction(client.NewGetSecretRequest(cfg, cmd.SecretIdentifier), cmd.SecretIdentifier)
	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}
