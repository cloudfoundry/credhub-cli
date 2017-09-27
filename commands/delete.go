package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type DeleteCommand struct {
	CredentialIdentifier string `short:"n" long:"name" required:"yes" description:"Name of the credential to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	credhubClient, err := initializeCredhubClient(cfg)

	if err != nil {
		return err
	}

	err = credhubClient.Delete(cmd.CredentialIdentifier)

	if err == nil {
		fmt.Println("Credential successfully deleted")
	}

	return err
}
