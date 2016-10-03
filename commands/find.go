package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type FindCommand struct {
	PartialSecretIdentifier string `short:"n" long:"name-like" description:"Find credentials whose name contains the query string"`
	PathIdentifier          string `short:"p" long:"path" description:"Find credentials that exist under the provided path"`
	AllPaths                bool   `short:"a" long:"all-paths" description:"List all existing credential paths"`
}

func (cmd FindCommand) Execute([]string) error {
	var credentials interface{}
	var err error
	var repository repositories.Repository

	cfg := config.ReadConfig()

	if cmd.AllPaths {
		repository = repositories.NewAllPathRepository(client.NewHttpClient(cfg))
	} else {
		repository = repositories.NewSecretQueryRepository(client.NewHttpClient(cfg))
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

	fmt.Println(credentials)

	return nil
}
