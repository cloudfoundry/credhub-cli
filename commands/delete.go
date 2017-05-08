package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type DeleteCommand struct {
	SecretIdentifier string `short:"n" long:"name" required:"yes" description:"Name of the credential to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, cfg)

	_, err := action.DoAction(client.NewDeleteSecretRequest(cfg, cmd.SecretIdentifier), cmd.SecretIdentifier)

	if err == nil {
		fmt.Println("Secret successfully deleted")
	}

	return err
}
