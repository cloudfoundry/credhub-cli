package commands

import (
	"fmt"

	"net/http"

	"os"

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
	request    *http.Request
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
		switch credential.Type {
		case "password", "value":
			value, ok := credential.Value.(string)
			if !ok {
				fmt.Fprintf(os.Stderr, "%v\n", "Interface conversion error")
				continue
			}
			request = client.NewSetCredentialRequest(cfg, credential.Type, credential.Name, value, true)
		case "certificate":
			certificate := new(models.Certificate)
			err = mapstructure.Decode(credential.Value, &certificate)

			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			request = client.NewSetCertificateRequest(cfg, credential.Name, certificate.Ca, certificate.CaName, certificate.Certificate, certificate.PrivateKey, true)
		case "rsa", "ssh":
			rsaSsh := new(models.RsaSsh)
			err = mapstructure.Decode(credential.Value, &rsaSsh)

			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			request = client.NewSetRsaSshRequest(cfg, credential.Name, credential.Type, rsaSsh.PublicKey, rsaSsh.PrivateKey, true)
		default:
			fmt.Fprintf(os.Stderr, "unrecognized type: %s", credential.Type)
		}

		result, err := action.DoAction(request, credential.Name)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		models.Println(result, false)
	}
}
