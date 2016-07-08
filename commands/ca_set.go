package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"

	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
)

type CaSetCommand struct {
	CaIdentifier      string `short:"n" required:"yes" long:"name" description:"Sets the name of the CA"`
	CaType            string `short:"t" long:"type" description:"Sets the type of the CA"`
	CaPublicFileName  string `long:"public" description:"Sets the Public Key based on an input file"`
	CaPrivateFileName string `long:"private" description:"Sets the Private Key based on an input file"`
	CaPublic          string `long:"public-string" description:"Sets the public key to the parameter value"`
	CaPrivate         string `long:"private-string" description:"Sets the private key to the parameter value"`
}

func (cmd CaSetCommand) Execute([]string) error {
	var err error

	config := config.ReadConfig()
	caRepository := repositories.NewCaRepository(client.NewHttpClient(config))

	action := actions.NewAction(caRepository, config)

	if cmd.CaPublicFileName != "" {
		if cmd.CaPublic != "" {
			return cmcli_errors.NewCombinationOfParametersError()
		}
		cmd.CaPublic, err = ReadFile(cmd.CaPublicFileName)
		if err != nil {
			return err
		}
	}
	if cmd.CaPrivateFileName != "" {
		if cmd.CaPrivate != "" {
			return cmcli_errors.NewCombinationOfParametersError()
		}
		cmd.CaPrivate, err = ReadFile(cmd.CaPrivateFileName)
		if err != nil {
			return err
		}
	}
	if cmd.CaType == "" {
		cmd.CaType = "root"
	}

	ca, err := action.DoAction(
		client.NewPutCaRequest(
			config.ApiURL,
			cmd.CaIdentifier,
			cmd.CaType,
			cmd.CaPublic,
			cmd.CaPrivate), cmd.CaIdentifier)

	if err != nil {
		return err
	}

	fmt.Println(ca)

	return nil
}
