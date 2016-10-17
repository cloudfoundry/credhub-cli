package commands

import (
	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/repositories"
	"github.com/pivotal-cf/credhub-cli/models"
)

type RegenerateCommand struct {
	SecretIdentifier string `required:"yes" short:"n" long:"name" description:"Selects the credential to regenerate"`
}

func (cmd RegenerateCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	repository := repositories.NewSecretRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, cfg)

	secret, err := action.DoAction(client.NewRegenerateSecretRequest(cfg, cmd.SecretIdentifier), cmd.SecretIdentifier)
	if err != nil {
		return err
	}

	models.Println(secret, false)

	return nil
}
