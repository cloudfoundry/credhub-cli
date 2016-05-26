package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type GetCommand struct {
	SecretIdentifier string `short:"n" required:"yes" long:"name" description:"Selects the secret to retrieve"`
}

func (cmd GetCommand) Execute([]string) error {
	action := actions.NewGet(client.NewHttpClient(), config.ReadConfig())

	secret, err := action.GetSecret(cmd.SecretIdentifier)
	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}
