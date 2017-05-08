package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type FindCommand struct {
	PartialSecretIdentifier string `short:"n" long:"name-like" description:"Find credentials whose name contains the query string"`
	PathIdentifier          string `short:"p" long:"path" description:"Find credentials that exist under the provided path"`
	AllPaths                bool   `short:"a" long:"all-paths" description:"List all existing credential paths"`
}

func (cmd FindCommand) Execute([]string) error {
	var credentials models.Printable
	var err error
	var repository repositories.Repository

	cfg := config.ReadConfig()

	if cmd.AllPaths {
		repository = repositories.NewAllPathRepository(client.NewHttpClient(cfg))
	} else {
		repository = repositories.NewCredentialQueryRepository(client.NewHttpClient(cfg))
	}

	action := actions.NewAction(repository, cfg)

	if cmd.AllPaths {
		credentials, err = action.DoAction(client.NewFindAllCredentialPathsRequest(cfg), "")
	} else if cmd.PartialSecretIdentifier != "" {
		credentials, err = action.DoAction(client.NewFindCredentialsBySubstringRequest(cfg, cmd.PartialSecretIdentifier), cmd.PartialSecretIdentifier)
	} else {
		credentials, err = action.DoAction(client.NewFindCredentialsByPathRequest(cfg, cmd.PathIdentifier), cmd.PartialSecretIdentifier)
	}
	if err != nil {
		return err
	}

	fmt.Println(credentials.ToYaml())

	return nil
}
