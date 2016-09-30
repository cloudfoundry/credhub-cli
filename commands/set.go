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
	SecretIdentifier  string `short:"n" required:"yes" long:"name" description:"Name of the credential to set"`
	Type              string `short:"t" long:"type" description:"Sets the credential type (Default: 'password')"`
	Value             string `short:"v" long:"value" description:"[Password, Value] Sets the value for the credential"`
	Root              string `short:"r" long:"root" description:"[Certificate] Sets the root CA from file"`
	Certificate       string `short:"c" long:"certificate" description:"[Certificate] Sets the certificate from file"`
	Private           string `short:"p" long:"private" description:"[Certificate, SSH] Sets the private key from file"`
	Public            string `short:"u" long:"public" description:"[SSH] Sets the public key from file"`
	RootString        string `short:"R" long:"root-string" description:"[Certificate] Sets the root CA from string input"`
	CertificateString string `short:"C" long:"certificate-string" description:"[Certificate] Sets the certificate from string input"`
	PrivateString     string `short:"P" long:"private-string" description:"[Certificate, SSH] Sets the private key from string input"`
	PublicString      string `short:"U" long:"public-string" description:"[SSH] Sets the public key from  string input"`
	NoOverwrite       bool   `short:"O" long:"no-overwrite" description:"Credential is not modified if stored value already exists"`
}

func (cmd SetCommand) Execute([]string) error {
	if cmd.Type == "" {
		cmd.Type = "password"
	}

	if cmd.Value == "" && (cmd.Type == "password" || cmd.Type == "value") {
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
	if cmd.Type == "value" {
		request = client.NewPutValueRequest(config, cmd.SecretIdentifier, cmd.Value, !cmd.NoOverwrite)
	} else if cmd.Type == "password" {
		request = client.NewPutPasswordRequest(config, cmd.SecretIdentifier, cmd.Value, !cmd.NoOverwrite)
	} else if cmd.Type == "ssh" {
		var err error

		err = setStringFieldFromFile(&cmd.Public, &cmd.PublicString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&cmd.Private, &cmd.PrivateString)
		if err != nil {
			return nil, err
		}

		request = client.NewPutSshRequest(config, cmd.SecretIdentifier, cmd.PublicString, cmd.PrivateString, !cmd.NoOverwrite)
	} else {
		var err error

		err = setStringFieldFromFile(&cmd.Root, &cmd.RootString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&cmd.Certificate, &cmd.CertificateString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&cmd.Private, &cmd.PrivateString)
		if err != nil {
			return nil, err
		}

		request = client.NewPutCertificateRequest(config, cmd.SecretIdentifier, cmd.RootString, cmd.CertificateString, cmd.PrivateString, !cmd.NoOverwrite)
	}

	return request, nil
}

func promptForInput(prompt string, value *string) {
	fmt.Printf(prompt)
	reader := bufio.NewReader(os.Stdin)
	val, _ := reader.ReadString('\n')
	*value = string(strings.TrimSpace(val))
}

func setStringFieldFromFile(fileField, stringField *string) error {
	var err error
	if *fileField != "" {
		if *stringField != "" {
			return cmcli_errors.NewCombinationOfParametersError()
		}
		*stringField, err = ReadFile(*fileField)
		if err != nil {
			return err
		}
	}
	return nil
}
