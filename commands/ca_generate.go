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
	CaCommonName       string `long:"common-name" description:"Sets the common name of the generated CA"`
	CaOrganization     string `long:"organization" description:"Sets the organization of the generated CA"`
	CaOrganizationUnit string `long:"organization-unit" description:"Sets the organization unit of the generated CA"`
	CaLocality         string `long:"locality" description:"Sets the locality/city of the generated CA"`
	CaState            string `long:"state" description:"Sets the state/province of the generated CA"`
	CaCountry          string `long:"country" description:"Sets the country of the generated CA"`
	CaKeyLength        int    `long:"key-length" description:"Sets the bit length of the generated CA key"`
	CaDuration         int    `long:"duration" description:"Sets the valid duration (in days) for the generated CA"`
}

func (cmd CaGenerateCommand) Execute([]string) error {
	var err error

	config, _ := config.ReadConfig()
	caRepository := repositories.NewCaRepository(client.NewHttpClient(config.ApiURL))

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

	action := actions.NewAction(caRepository, config)

	if cmd.CaType == "" {
		cmd.CaType = "root"
	}

	request := client.NewPostCaRequest(config, cmd.CaIdentifier, cmd.CaType, parameters)

	ca, err := action.DoAction(request, cmd.CaIdentifier)

	if err != nil {
		return err
	}

	fmt.Println(ca)

	return nil
}
