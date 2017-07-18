package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type GetCommand struct {
	Name       string `short:"n" long:"name" description:"Name of the credential to retrieve"`
	Id         string `long:"id" description:"ID of the credential to retrieve"`
	OutputJson bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd GetCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	a := api.NewApi(&cfg)

	credential, err := a.Get(cmd.Name, cmd.Id)

	if err != nil {
		return err
	}

	models.Println(credential, cmd.OutputJson)

	return nil
}
