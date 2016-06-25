package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type CaSetCommand struct {
	CaIdentifier string `short:"n" required:"yes" long:"name" description:"Sets the name of the CA"`
	CaPublic     string `long:"public-string" description:"Sets the public key to the parameter value"`
	CaPrivate    string `long:"private-string" description:"Sets the private key to the parameter value"`
}

func (cmd CaSetCommand) Execute([]string) error {
	caRepository := repositories.NewCaRepository(client.NewHttpClient())

	config := config.ReadConfig()
	action := actions.NewAction(caRepository, config)
	ca, err := action.DoAction(
		client.NewPutCaRequest(
			config.ApiURL,
			cmd.CaIdentifier,
			cmd.CaPublic,
			cmd.CaPrivate), cmd.CaIdentifier)

	if err != nil {
		return err
	}

	fmt.Println(ca)

	return nil
}
