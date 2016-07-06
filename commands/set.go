package commands

import (
	"fmt"

	"net/http"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type SetCommand struct {
	SecretIdentifier           string `short:"n" required:"yes" long:"name" description:"Selects the secret being set"`
	ContentType                string `short:"t" long:"type" description:"Sets the type of secret to store or generate. Default: 'value'"`
	Value                      string `short:"v" long:"value" description:"Sets a value for a secret name"`
	CertificateCAFileName      string `long:"ca" description:"Sets the CA based on an input file"`
	CertificatePublicFileName  string `long:"public" description:"Sets the Public Key based on an input file"`
	CertificatePrivateFileName string `long:"private" description:"Sets the Private Key based on an input file"`
	CertificateCA              string `long:"ca-string" description:"Sets the Certificate Authority"`
	CertificatePublic          string `long:"public-string" description:"Sets the Public Key"`
	CertificatePrivate         string `long:"private-string" description:"Sets the Private Key"`
}

func (cmd SetCommand) Execute([]string) error {
	if cmd.ContentType == "" {
		cmd.ContentType = "value"
	}

	repository := repositories.NewSecretRepository(client.NewHttpClient())

	config := config.ReadConfig()
	action := actions.NewAction(repository, config)
	request, err := getRequest(cmd, config.ApiURL)
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

func getRequest(cmd SetCommand, url string) (*http.Request, error) {
	var request *http.Request
	if cmd.ContentType == "value" {
		request = client.NewPutValueRequest(url, cmd.SecretIdentifier, cmd.Value)
	} else {
		var err error
		if cmd.CertificateCAFileName != "" {
			cmd.CertificateCA, err = ReadFile(cmd.CertificateCAFileName)
			if err != nil {
				return nil, err
			}
		}
		if cmd.CertificatePublicFileName != "" {
			cmd.CertificatePublic, err = ReadFile(cmd.CertificatePublicFileName)
			if err != nil {
				return nil, err
			}
		}
		if cmd.CertificatePrivateFileName != "" {
			cmd.CertificatePrivate, err = ReadFile(cmd.CertificatePrivateFileName)
			if err != nil {
				return nil, err
			}
		}
		request = client.NewPutCertificateRequest(url, cmd.SecretIdentifier, cmd.CertificateCA, cmd.CertificatePublic, cmd.CertificatePrivate)
	}

	return request, nil
}
