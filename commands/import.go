package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
	"github.com/mitchellh/mapstructure"
)

type ImportCommand struct {
	File string `short:"f" long:"file" description:"File containing credentials to import"`
}

var (
	err        error
	repository repositories.Repository
	bulkImport models.CredentialBulkImport
	setCommand SetCommand
)

func (cmd ImportCommand) Execute([]string) error {

	err = bulkImport.ReadFile(cmd.File)

	if err != nil {
		return err
	}

	setCredentials(bulkImport)

	return nil
}

func setCredentials(bulkImport models.CredentialBulkImport) {
	cfg := config.ReadConfig()
	repository = repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, cfg)

	for _, credential := range bulkImport.Credentials {
		setCommand = SetCommand{}
		setCommand.CredentialIdentifier = credential.Name
		setCommand.Type = credential.Type

		switch credential.Type {
		case "password":
			setCommand.Password = credential.Value.(string)
		case "value":
			setCommand.Value = credential.Value.(string)
		case "certificate":
			certificate := new(models.Certificate)
			err = mapstructure.Decode(credential.Value, &certificate)

			if certificate.CaName != "" {
				setCommand.CaName = certificate.CaName
			}

			if certificate.Ca != "" {
				setCommand.RootString = certificate.Ca
			}

			setCommand.CertificateString = certificate.Certificate
			setCommand.PrivateString = certificate.PrivateKey
		default:
			fmt.Errorf("unrecognized type: %s", credential.Type)
		}

		setRequest, err := MakeRequest(setCommand, cfg)

		if err != nil {
			fmt.Errorf("%v\n", err)
			continue
		}

		result, err := action.DoAction(setRequest, credential.Name)

		models.Println(result, false)

		if err != nil {
			fmt.Errorf("%v\n", err)
			continue
		}
	}
}
