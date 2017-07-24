package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type DeleteCommand struct {
	CredentialIdentifier string `short:"n" long:"name" required:"yes" description:"Name of the credential to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	err := api.NewApi(&cfg).Delete(cmd.CredentialIdentifier)

	if err == nil {
		fmt.Println("Credential successfully deleted")
	}

	return err
}
