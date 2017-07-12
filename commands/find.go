package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type FindCommand struct {
	PartialCredentialIdentifier string `short:"n" long:"name-like" description:"Find credentials whose name contains the query string"`
	PathIdentifier              string `short:"p" long:"path" description:"Find credentials that exist under the provided path"`
	AllPaths                    bool   `short:"a" long:"all-paths" description:"List all existing credential paths"`
	OutputJson                  bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd FindCommand) Execute([]string) error {
	credentials, err := api.Find(cmd.PartialCredentialIdentifier, cmd.PathIdentifier, cmd.AllPaths)

	if err != nil {
		return err
	}

	models.Println(credentials, cmd.OutputJson)

	return nil
}
