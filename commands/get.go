package commands

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/config"
)

type GetCommand struct {
	SecretIdentifier string `short:"n" required:"yes" long:"name" description:"Selects the secret to retrieve"`
}

func (cmd GetCommand) Execute([]string) error {
	action := actions.NewGet(http.DefaultClient, config.ReadConfig())

	secret, err := action.GetSecret(cmd.SecretIdentifier)
	if err != nil {
		return err
	}

	secret.PrintSecret()

	return nil
}
