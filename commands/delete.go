package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type DeleteCommand struct {
	SecretIdentifier string `short:"n" long:"name" required:"yes" description:"Selects the secret to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	action := actions.NewDelete(client.NewHttpClient(), config.ReadConfig())

	err := action.Delete(cmd.SecretIdentifier)

	if err == nil {
		fmt.Println("Secret successfully deleted")
	}

	return err
}
