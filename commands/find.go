package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/models"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type FindCommand struct {
	PartialSecretIdentifier string `short:"n" long:"name-like" description:"Find credentials whose name contains the query string"`
	PathIdentifier          string `short:"p" long:"path" description:"Find credentials that exist under the provided path"`
}

func (cmd FindCommand) Execute([]string) error {
	var credentials models.Item
	var err error

	cfg := config.ReadConfig()
	repository := repositories.NewSecretQueryRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, cfg)

	if cmd.PartialSecretIdentifier != "" {
		credentials, err = action.DoAction(client.NewFindCredentialsBySubstringRequest(cfg, cmd.PartialSecretIdentifier), cmd.PartialSecretIdentifier)
	} else {
		credentials, err = action.DoAction(client.NewFindCredentialsByPathRequest(cfg, cmd.PathIdentifier), cmd.PartialSecretIdentifier)
	}
	if err != nil {
		return err
	}

	fmt.Println(credentials)

	return nil
}
