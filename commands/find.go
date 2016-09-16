package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type FindCommand struct {
	PartialSecretIdentifier string `short:"n" required:"yes" long:"name-like" description:"Find credentials whose name contains the query string"`
}

func (cmd FindCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	repository := repositories.NewSecretQueryRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, cfg)
	secret, err := action.DoAction(client.NewFindSecretsRequest(cfg, cmd.PartialSecretIdentifier), cmd.PartialSecretIdentifier)
	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}
