package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/config"

	"bufio"
	"os"
	"strings"

	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/util"
)

type SetCommand struct {
	CredentialIdentifier string `short:"n" required:"yes" long:"name" description:"Name of the credential to set"`
	Type                 string `short:"t" long:"type" description:"Sets the credential type. Valid types include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'."`
	NoOverwrite          bool   `short:"O" long:"no-overwrite" description:"Credential is not modified if stored value already exists"`
	Value                string `short:"v" long:"value" description:"[Value, JSON] Sets the value for the credential"`
	CaName               string `short:"m" long:"ca-name" description:"[Certificate] Sets the root CA to a stored CA credential"`
	Root                 string `short:"r" long:"root" description:"[Certificate] Sets the root CA from file or value"`
	Certificate          string `short:"c" long:"certificate" description:"[Certificate] Sets the certificate from file or value"`
	Private              string `short:"p" long:"private" description:"[Certificate, SSH, RSA] Sets the private key from file or value"`
	Public               string `short:"u" long:"public" description:"[SSH, RSA] Sets the public key from file or value"`
	Username             string `short:"z" long:"username" description:"[User] Sets the username value of the credential"`
	Password             string `short:"w" long:"password" description:"[Password, User] Sets the password value of the credential"`
	OutputJson           bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
}

func (cmd SetCommand) Execute([]string) error {
	cmd.Type = strings.ToLower(cmd.Type)

	if cmd.Type == "" {
		return errors.NewSetEmptyTypeError()
	}

	cmd.setFieldsFromInteractiveUserInput()

	err := cmd.setFieldsFromFileOrString()
	if err != nil {
		return err
	}

	cfg := config.ReadConfig()

	credhubClient, err := initializeCredhubClient(cfg)
	if err != nil {
		return err
	}

	credential, err := cmd.setCredential(credhubClient)
	if err != nil {
		return err
	}

	printCredential(cmd.OutputJson, credential)

	return nil
}

func (cmd *SetCommand) setFieldsFromInteractiveUserInput() {
	if cmd.Value == "" && (cmd.Type == "value" || cmd.Type == "json") {
		promptForInput("value: ", &cmd.Value)
	}

	if cmd.Password == "" && (cmd.Type == "password" || cmd.Type == "user") {
		promptForInput("password: ", &cmd.Password)
	}
}

func (cmd *SetCommand) setFieldsFromFileOrString() error {
	var err error

	cmd.Public, err = util.ReadFileOrStringFromField(cmd.Public)
	if err != nil {
		return err
	}

	cmd.Private, err = util.ReadFileOrStringFromField(cmd.Private)
	if err != nil {
		return err
	}

	cmd.Root, err = util.ReadFileOrStringFromField(cmd.Root)
	if err != nil {
		return err
	}

	cmd.Certificate, err = util.ReadFileOrStringFromField(cmd.Certificate)

	return err
}

func (cmd SetCommand) setCredential(credhubClient *credhub.CredHub) (interface{}, error) {
	mode := credhub.Overwrite

	if cmd.NoOverwrite {
		mode = credhub.NoOverwrite
	}

	var value interface{}

	switch cmd.Type {
	case "password":
		value = values.Password(cmd.Password)
	case "certificate":
		value = values.Certificate{
			Ca:          cmd.Root,
			Certificate: cmd.Certificate,
			PrivateKey:  cmd.Private,
			CaName:      cmd.CaName,
		}
	case "ssh":
		value = values.SSH{
			PublicKey:  cmd.Public,
			PrivateKey: cmd.Private,
		}
	case "rsa":
		value = values.RSA{
			PublicKey:  cmd.Public,
			PrivateKey: cmd.Private,
		}
	case "user":
		value = values.User{
			Password: cmd.Password,
			Username: cmd.Username,
		}
	case "json":
		value = values.JSON{}
		err := json.Unmarshal([]byte(cmd.Value), &value)
		if err != nil {
			return nil, err
		}
	default:
		value = values.Value(cmd.Value)
	}
	return credhubClient.SetCredential(cmd.CredentialIdentifier, cmd.Type, value, mode)
}

func promptForInput(prompt string, value *string) {
	fmt.Printf(prompt)
	reader := bufio.NewReader(os.Stdin)
	val, _ := reader.ReadString('\n')
	*value = string(strings.TrimSpace(val))
}
