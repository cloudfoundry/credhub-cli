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
	CaIdentifier      string `short:"n" required:"yes" long:"name" description:"Name of the CA to set"`
	CaType            string `short:"t" long:"type" description:"Sets the CA type (Default: 'root')"`
	CaPublicFileName  string `short:"c" long:"certificate" description:"[Certificate] Sets the CA certificate from file"`
	CaPrivateFileName string `short:"p" long:"private" description:"[Certificate] Sets the CA private key from file"`
	CaPublic          string `short:"C" long:"certificate-string" description:"[Certificate] Sets the CA certificate from string input"`
	CaPrivate         string `short:"P" long:"private-string" description:"[Certificate] Sets the CA private key from string input"`
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

	fmt.Println(ca.Terminal())

	return nil
}
