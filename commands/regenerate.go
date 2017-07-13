package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type RegenerateCommand struct {
	CredentialIdentifier string `required:"yes" short:"n" long:"name" description:"Selects the credential to regenerate"`
	OutputJson           bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd RegenerateCommand) Execute([]string) error {
	credential, err := api.Regenerate(cmd.CredentialIdentifier)

	if err != nil {
		return err
	}

	models.Println(credential, cmd.OutputJson)
	return nil
}
