package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/repositories"

	cmcli_errors "github.com/pivotal-cf/credhub-cli/errors"
)

type CaSetCommand struct {
	CaIdentifier      string `short:"n" required:"yes" long:"name" description:"Sets the name of the CA"`
	CaType            string `short:"t" long:"type" description:"Sets the type of the CA"`
	CaPublicFileName  string `short:"c" long:"certificate" description:"Sets the Certificate based on an input file"`
	CaPrivateFileName string `short:"p" long:"private" description:"Sets the Private Key based on an input file"`
	CaPublic          string `short:"C" long:"certificate-string" description:"Sets the Certificate to the parameter value"`
	CaPrivate         string `short:"P" long:"private-string" description:"Sets the Private Key to the parameter value"`
}

func (cmd CaSetCommand) Execute([]string) error {
	var err error

	cfg := config.ReadConfig()
	caRepository := repositories.NewCaRepository(client.NewHttpClient(cfg))

	action := actions.NewAction(caRepository, cfg)

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
			cfg,
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
