package commands

import (
	"fmt"

	"net/http"

	"io/ioutil"

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

	secretRepository := repositories.NewSecretRepository(client.NewHttpClient())

	config := config.ReadConfig()
	action := actions.NewAction(secretRepository, config)
	secret, err := action.DoAction(getRequest(cmd, config.ApiURL), cmd.SecretIdentifier)

	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}

func getRequest(cmd SetCommand, url string) *http.Request {
	var request *http.Request
	if cmd.ContentType == "value" {
		request = client.NewPutValueRequest(url, cmd.SecretIdentifier, cmd.Value)
	} else {
		if cmd.CertificateCAFileName != "" {
			cmd.CertificateCA = readFile(cmd.CertificateCAFileName)
		}
		if cmd.CertificatePublicFileName != "" {
			cmd.CertificatePublic = readFile(cmd.CertificatePublicFileName)
		}
		if cmd.CertificatePrivateFileName != "" {
			cmd.CertificatePrivate = readFile(cmd.CertificatePrivateFileName)
		}
		request = client.NewPutCertificateRequest(url, cmd.SecretIdentifier, cmd.CertificateCA, cmd.CertificatePublic, cmd.CertificatePrivate)
	}

	return request
}

func readFile(filename string) string {
	dat, _ := ioutil.ReadFile(filename)
	return string(dat)
}
