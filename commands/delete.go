package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
)

type DeleteCommand struct {
	CredentialIdentifier string `short:"n" long:"name" required:"yes" description:"Name of the credential to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	err := api.Delete(cmd.CredentialIdentifier)

	if err == nil {
		fmt.Println("Credential successfully deleted")
	}

	return err
}
