package commands

import (
	"fmt"

	"net/http"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"

	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
)

type SetCommand struct {
	SecretIdentifier           string `short:"n" required:"yes" long:"name" description:"Selects the secret being set"`
	ContentType                string `short:"t" long:"type" description:"Sets the type of secret to store or generate. Default: 'value'"`
	Value                      string `short:"v" long:"value" description:"Sets a value for a secret name"`
	RootCAFileName             string `long:"root" description:"Sets the Root CA based on an input file"`
	CertificatePublicFileName  string `long:"certificate" description:"Sets the Certificate based on an input file"`
	CertificatePrivateFileName string `long:"private" description:"Sets the Private Key based on an input file"`
	RootCA                     string `long:"root-string" description:"Sets the Root Certificate Authority"`
	CertificatePublic          string `long:"certificate-string" description:"Sets the Certificate"`
	CertificatePrivate         string `long:"private-string" description:"Sets the Private Key"`
}

func (cmd SetCommand) Execute([]string) error {
	if cmd.ContentType == "" {
		cmd.ContentType = "value"
	}

	config, _ := config.ReadConfig()
	repository := repositories.NewSecretRepository(client.NewHttpClient(config.ApiURL))

	action := actions.NewAction(repository, config)
	request, err := getRequest(cmd, config)
	if err != nil {
		return err
	}

	secret, err := action.DoAction(request, cmd.SecretIdentifier)
	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}

func getRequest(cmd SetCommand, config config.Config) (*http.Request, error) {
	var request *http.Request
	if cmd.ContentType == "value" {
		request = client.NewPutValueRequest(config, cmd.SecretIdentifier, cmd.Value)
	} else {
		var err error
		if cmd.RootCAFileName != "" {
			if cmd.RootCA != "" {
				return nil, cmcli_errors.NewCombinationOfParametersError()
			}
			cmd.RootCA, err = ReadFile(cmd.RootCAFileName)
			if err != nil {
				return nil, err
			}
		}
		if cmd.CertificatePublicFileName != "" {
			if cmd.CertificatePublic != "" {
				return nil, cmcli_errors.NewCombinationOfParametersError()
			}
			cmd.CertificatePublic, err = ReadFile(cmd.CertificatePublicFileName)
			if err != nil {
				return nil, err
			}
		}
		if cmd.CertificatePrivateFileName != "" {
			if cmd.CertificatePrivate != "" {
				return nil, cmcli_errors.NewCombinationOfParametersError()
			}
			cmd.CertificatePrivate, err = ReadFile(cmd.CertificatePrivateFileName)
			if err != nil {
				return nil, err
			}
		}
		request = client.NewPutCertificateRequest(config, cmd.SecretIdentifier, cmd.RootCA, cmd.CertificatePublic, cmd.CertificatePrivate)
	}

	return request, nil
}
