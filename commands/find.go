package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
)

type FindCommand struct {
	PartialCredentialIdentifier string `short:"n" long:"name-like" description:"Find credentials whose name contains the query string"`
	PathIdentifier              string `short:"p" long:"path" description:"Find credentials that exist under the provided path"`
	AllPaths                    bool   `short:"a" long:"all-paths" description:"List all existing credential paths"`
	OutputJson                  bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd FindCommand) Execute([]string) error {
	var output interface{}
	var err error

	cfg := config.ReadConfig()

	credhubClient, err := initializeCredhubClient(cfg)
	if err != nil {
		return err
	}

	if cmd.AllPaths {
		var paths credentials.Paths
		paths, err = credhubClient.FindAllPaths()
		if len(paths.Paths) == 0 {
			return errors.NewNoMatchingCredentialsFoundError()
		}
		output = paths
	} else if cmd.PartialCredentialIdentifier != "" {
		var results credentials.FindResults
		results, err = credhubClient.FindByPartialName(cmd.PartialCredentialIdentifier)
		if len(results.Credentials) == 0 {
			return errors.NewNoMatchingCredentialsFoundError()
		}
		output = results
	} else {
		output, err = credhubClient.FindByPath(cmd.PathIdentifier)
	}

	if err != nil {
		return err
	}

	printCredential(cmd.OutputJson, output)

	return nil
}
