package commands

import (
	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/models"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type CaGenerateCommand struct {
	CaIdentifier       string `short:"n" required:"yes" long:"name" description:"Name of the CA to generate"`
	CaType             string `short:"t" long:"type" description:"Sets the CA type to generate (Default: 'root')"`
	CaDuration         int    `short:"d" long:"duration" description:"[Root] Valid duration (in days) of the generated CA certificate (Default: 365)"`
	CaKeyLength        int    `short:"k" long:"key-length" description:"[Root] Bit length of the generated key (Default: 2048)"`
	CaCommonName       string `short:"c" long:"common-name" description:"[Root] Common name of the generated CA certificate"`
	CaOrganization     string `short:"o" long:"organization" description:"[Root] Organization of the generated CA certificate"`
	CaOrganizationUnit string `short:"u" long:"organization-unit" description:"[Root] Organization unit of the generated CA certificate"`
	CaLocality         string `short:"i" long:"locality" description:"[Root] Locality/city of the generated CA certificate"`
	CaState            string `short:"s" long:"state" description:"[Root] State/province of the generated CA certificate"`
	CaCountry          string `short:"y" long:"country" description:"[Root] Country of the generated CA certificate"`
	OutputJson         bool   `long:"output-json"`
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

	models.Println(ca, cmd.OutputJson)

	return nil
}
