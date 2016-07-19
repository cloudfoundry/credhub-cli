package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type CaGenerateCommand struct {
	CaIdentifier string `short:"n" required:"yes" long:"name" description:"Sets the name of the CA"`
	CaType       string `short:"t" long:"type" description:"Sets the type of the CA"`
}

func (cmd CaGenerateCommand) Execute([]string) error {
	var err error

	config := config.ReadConfig()
	caRepository := repositories.NewCaRepository(client.NewHttpClient(config.ApiURL))

	action := actions.NewAction(caRepository, config)

	if cmd.CaType == "" {
		cmd.CaType = "root"
	}

	ca, err := action.DoAction(
		client.NewPostCaRequest(
			config.ApiURL,
			cmd.CaIdentifier,
			cmd.CaType), cmd.CaIdentifier)

	if err != nil {
		return err
	}

	fmt.Println(ca)

	return nil
}
