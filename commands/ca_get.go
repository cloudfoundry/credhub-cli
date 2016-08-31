package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type CaGetCommand struct {
	CaIdentifier string `short:"n" required:"yes" long:"name" description:"Name of the CA to retrieve"`
}

func (cmd CaGetCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	caRepository := repositories.NewCaRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(caRepository, cfg)

	ca, err := action.DoAction(client.NewGetCaRequest(cfg, cmd.CaIdentifier), cmd.CaIdentifier)

	if err != nil {
		return err
	}

	fmt.Println(ca)

	return nil
}
