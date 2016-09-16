package commands

import (
	"fmt"

	"net/http"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/repositories"

	"bufio"
	"os"
	"strings"

	cmcli_errors "github.com/pivotal-cf/credhub-cli/errors"
)

type SetCommand struct {
	SecretIdentifier           string `short:"n" required:"yes" long:"name" description:"Name of the credential to set"`
	ContentType                string `short:"t" long:"type" description:"Sets the credential type (Default: 'password')"`
	Value                      string `short:"v" long:"value" description:"[Password, Value] Sets the value for the credential"`
	RootCAFileName             string `short:"r" long:"root" description:"[Certificate] Sets the root CA from file"`
	CertificatePublicFileName  string `short:"c" long:"certificate" description:"[Certificate] Sets the certificate from file"`
	CertificatePrivateFileName string `short:"p" long:"private" description:"[Certificate] Sets the private key from file"`
	RootCA                     string `short:"R" long:"root-string" description:"[Certificate] Sets the root CA from string input"`
	CertificatePublic          string `short:"C" long:"certificate-string" description:"[Certificate] Sets the certificate from string input"`
	CertificatePrivate         string `short:"P" long:"private-string" description:"[Certificate] Sets the private key from string input"`
	NoOverwrite                bool   `short:"O" long:"no-overwrite" description:"Credential is not modified if stored value already exists"`
}

func (cmd SetCommand) Execute([]string) error {
	if cmd.ContentType == "" {
		cmd.ContentType = "password"
	}

	if cmd.Value == "" && (cmd.ContentType == "password" || cmd.ContentType == "value") {
		promptForInput("value: ", &cmd.Value)
	}

	cfg := config.ReadConfig()
	repository := repositories.NewSecretRepository(client.NewHttpClient(cfg))

	action := actions.NewAction(repository, cfg)
	request, err := getRequest(cmd, cfg)
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
		request = client.NewPutValueRequest(config, cmd.SecretIdentifier, cmd.Value, !cmd.NoOverwrite)
	} else if cmd.ContentType == "password" {
		request = client.NewPutPasswordRequest(config, cmd.SecretIdentifier, cmd.Value, !cmd.NoOverwrite)
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
		request = client.NewPutCertificateRequest(config, cmd.SecretIdentifier, cmd.RootCA, cmd.CertificatePublic, cmd.CertificatePrivate, !cmd.NoOverwrite)
	}

	return request, nil
}

func promptForInput(prompt string, value *string) {
	fmt.Printf(prompt)
	reader := bufio.NewReader(os.Stdin)
	val, _ := reader.ReadString('\n')
	*value = string(strings.TrimSpace(val))
}
