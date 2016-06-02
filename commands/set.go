package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type SetCommand struct {
	SecretIdentifier string `short:"n" required:"yes" long:"name" description:"Selects the secret being set"`
	ContentType      string `short:"t" long:"type" description:"Sets the type of secret to store or generate. Default: 'value'"`
	SecretContent    string `short:"v" required:"yes" long:"value" description:"Sets a value for a secret name"`
}

func (cmd SetCommand) Execute([]string) error {
	if cmd.SecretContent == "" {
		return errors.NewSetOptionMissingError()
	}

	if cmd.ContentType == "" {
		cmd.ContentType = "value"
	}

	secretRepository := repositories.NewSecretRepository(client.NewHttpClient())

	action := actions.NewSet(secretRepository, config.ReadConfig())
	secret, err := action.SetSecret(cmd.SecretIdentifier, cmd.SecretContent, cmd.ContentType)

	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}
