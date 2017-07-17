package commands

import (
	"fmt"

	"bufio"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type SetCommand struct {
	CredentialIdentifier string `short:"n" required:"yes" long:"name" description:"Name of the credential to set"`
	Type                 string `short:"t" long:"type" description:"Sets the credential type. Valid types include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'. Type-specific values are set with the following flags (supported types prefixed)."`
	NoOverwrite          bool   `short:"O" long:"no-overwrite" description:"Credential is not modified if stored value already exists"`
	Value                string `short:"v" long:"value" description:"[Value, JSON] Sets the value for the credential"`
	CaName               string `short:"m" long:"ca-name" description:"[Certificate] Sets the root CA to a stored CA credential"`
	Root                 string `short:"r" long:"root" description:"[Certificate] Sets the root CA from file"`
	Certificate          string `short:"c" long:"certificate" description:"[Certificate] Sets the certificate from file"`
	Private              string `short:"p" long:"private" description:"[Certificate, SSH, RSA] Sets the private key from file"`
	Public               string `short:"u" long:"public" description:"[SSH, RSA] Sets the public key from file"`
	RootString           string `short:"R" long:"root-string" description:"[Certificate] Sets the root CA from string input"`
	CertificateString    string `short:"C" long:"certificate-string" description:"[Certificate] Sets the certificate from string input"`
	PrivateString        string `short:"P" long:"private-string" description:"[Certificate, SSH, RSA] Sets the private key from string input"`
	PublicString         string `short:"U" long:"public-string" description:"[SSH, RSA] Sets the public key from  string input"`
	Username             string `short:"z" long:"username" description:"[User] Sets the username value of the credential"`
	Password             string `short:"w" long:"password" description:"[Password, User] Sets the password value of the credential"`
	OutputJson           bool   `          long:"output-json" description:"Return response in JSON format"`
}

func (cmd SetCommand) Execute([]string) error {

	if cmd.Type == "" {
		return errors.NewSetEmptyTypeError()
	}

	if cmd.Value == "" && (cmd.Type == "value" || cmd.Type == "json") {
		promptForInput("value: ", &cmd.Value)
	}

	if cmd.Password == "" && (cmd.Type == "password" || cmd.Type == "user") {
		promptForInput("password: ", &cmd.Password)
	}

	credential, err := api.Set(
		cmd.CredentialIdentifier,
		cmd.Type,
		cmd.NoOverwrite,
		cmd.Value,
		cmd.CaName,
		cmd.Root,
		cmd.Certificate,
		cmd.Private,
		cmd.Public,
		cmd.RootString,
		cmd.CertificateString,
		cmd.PrivateString,
		cmd.PublicString,
		cmd.Username,
		cmd.Password,
	)

	if err != nil {
		return err
	}

	models.Println(credential, cmd.OutputJson)

	return nil
}

func promptForInput(prompt string, value *string) {
	fmt.Printf(prompt)
	reader := bufio.NewReader(os.Stdin)
	val, _ := reader.ReadString('\n')
	*value = string(strings.TrimSpace(val))
}
