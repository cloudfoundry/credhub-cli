package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/models"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type CaGenerateCommand struct {
	CaIdentifier       string `short:"n" required:"yes" long:"name" description:"Sets the name of the CA"`
	CaType             string `short:"t" long:"type" description:"Sets the type of the CA"`
	CaCommonName       string `short:"c" long:"common-name" description:"Sets the common name of the generated CA"`
	CaOrganization     string `short:"o" long:"organization" description:"Sets the organization of the generated CA"`
	CaOrganizationUnit string `short:"u" long:"organization-unit" description:"Sets the organization unit of the generated CA"`
	CaLocality         string `short:"i" long:"locality" description:"Sets the locality/city of the generated CA"`
	CaState            string `short:"s" long:"state" description:"Sets the state/province of the generated CA"`
	CaCountry          string `short:"y" long:"country" description:"Sets the country of the generated CA"`
	CaKeyLength        int    `short:"k" long:"key-length" description:"Sets the bit length of the generated CA key"`
	CaDuration         int    `short:"d" long:"duration" description:"Sets the valid duration (in days) for the generated CA"`
}

func (cmd CaGenerateCommand) Execute([]string) error {
	var err error

	cfg := config.ReadConfig()
	caRepository := repositories.NewCaRepository(client.NewHttpClient(cfg))

	parameters := models.SecretParameters{
		CommonName:       cmd.CaCommonName,
		Organization:     cmd.CaOrganization,
		OrganizationUnit: cmd.CaOrganizationUnit,
		Locality:         cmd.CaLocality,
		State:            cmd.CaState,
		Country:          cmd.CaCountry,
		KeyLength:        cmd.CaKeyLength,
		Duration:         cmd.CaDuration,
	}

	action := actions.NewAction(caRepository, cfg)

	if cmd.CaType == "" {
		cmd.CaType = "root"
	}

	request := client.NewPostCaRequest(cfg, cmd.CaIdentifier, cmd.CaType, parameters)

	ca, err := action.DoAction(request, cmd.CaIdentifier)

	if err != nil {
		return err
	}

	fmt.Println(ca)

	return nil
}
