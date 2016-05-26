package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type SetCommand struct {
	SecretIdentifier string `short:"n" required:"yes" long:"name" description:"Selects the secret being set"`
	SecretContent    string `short:"s" long:"secret" description:"Sets a value for a secret name"`
	Generate         bool   `short:"g" long:"generate" description:"System will generate random credential. Cannot be used in combination with --secret."`
	Length           int    `short:"l" long:"length" description:"Sets length of generated value (Default: 20)"`
}

func (cmd SetCommand) Execute([]string) error {
	if !cmd.Generate && cmd.SecretContent == "" {
		return errors.NewSetOptionMissingError()
	}

	secretRepository := repositories.NewSecretRepository(client.NewHttpClient())

	var secret models.Secret
	var err error

	if cmd.Generate {
		parameters := models.SecretParameters{
			Length: cmd.Length,
		}

		action := actions.NewGenerate(secretRepository, config.ReadConfig())
		secret, err = action.GenerateSecret(cmd.SecretIdentifier, parameters)
	} else {
		action := actions.NewSet(secretRepository, config.ReadConfig())
		secret, err = action.SetSecret(cmd.SecretIdentifier, cmd.SecretContent)
	}

	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}
