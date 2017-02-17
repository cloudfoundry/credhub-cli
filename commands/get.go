package commands

import (
	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/models"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type GetCommand struct {
	SecretIdentifier string `short:"n" required:"yes" long:"name" description:"Name of the credential to retrieve"`
	OutputJson       bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd GetCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	repository := repositories.NewSecretRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, cfg)
	secret, err := action.DoAction(client.NewGetSecretRequest(cfg, cmd.SecretIdentifier), cmd.SecretIdentifier)
	if err != nil {
		return err
	}

	models.Println(secret, cmd.OutputJson)

	return nil
}
