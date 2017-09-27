package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/config"

	"bufio"
	"os"
	"strings"

	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
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
	OutputJson           bool   `          long:"output-json" description:"Return response in JSON format"`
}

func (cmd SetCommand) Execute([]string) error {
	cmd.Type = strings.ToLower(cmd.Type)

	if cmd.Type == "" {
		return errors.NewSetEmptyTypeError()
	}

	if cmd.Value == "" && (cmd.Type == "value" || cmd.Type == "json") {
		promptForInput("value: ", &cmd.Value)
	}

	if cmd.Password == "" && (cmd.Type == "password" || cmd.Type == "user") {
		promptForInput("password: ", &cmd.Password)
	}

	cfg := config.ReadConfig()

	credhubClient, err := initializeCredhubClient(cfg)
	if err != nil {
		return err
	}

	err = config.ValidateConfig(cfg)
	if err != nil {
		if !clientCredentialsInEnvironment() || config.ValidateConfigApi(cfg) != nil {
			return err
		}
	}

	credential, err := MakeRequest(cmd, cfg, credhubClient)
	if err != nil {
		return err
	}

	printCredential(cmd.OutputJson, credential)

	return nil
}

func MakeRequest(cmd SetCommand, config config.Config, credhubClient *credhub.CredHub) (interface{}, error) {
	var output interface{}
	var responseError error

	if cmd.Type == "ssh" || cmd.Type == "rsa" {
		publicKey, err := util.ReadFileOrStringFromField(cmd.Public)
		if err != nil {
			return nil, err
		}

		privateKey, err := util.ReadFileOrStringFromField(cmd.Private)
		if err != nil {
			return nil, err
		}
		if cmd.Type == "ssh" {
			value := values.SSH{}
			value.PublicKey = publicKey
			value.PrivateKey = privateKey
			var sshCredential credentials.SSH
			sshCredential, responseError = credhubClient.SetSSH(cmd.CredentialIdentifier, value, !cmd.NoOverwrite)
			output = interface{}(sshCredential)
		} else {
			value := values.RSA{}
			value.PublicKey = publicKey
			value.PrivateKey = privateKey
			var rsaCredential credentials.RSA
			rsaCredential, responseError = credhubClient.SetRSA(cmd.CredentialIdentifier, value, !cmd.NoOverwrite)
			output = interface{}(rsaCredential)
		}
	} else if cmd.Type == "certificate" {

		root, err := util.ReadFileOrStringFromField(cmd.Root)
		if err != nil {
			return nil, err
		}

		certificate, err := util.ReadFileOrStringFromField(cmd.Certificate)
		if err != nil {
			return nil, err
		}

		privateKey, err := util.ReadFileOrStringFromField(cmd.Private)
		if err != nil {
			return nil, err
		}

		value := values.Certificate{}
		value.Certificate = certificate
		value.PrivateKey = privateKey
		value.Ca = root
		value.CaName = cmd.CaName
		var certificateCredential credentials.Certificate
		certificateCredential, responseError = credhubClient.SetCertificate(cmd.CredentialIdentifier, value, !cmd.NoOverwrite)
		output = interface{}(certificateCredential)
	} else if cmd.Type == "user" {
		value := values.User{}
		if cmd.Username != "" {
			value.Username = &cmd.Username
		}
		value.Password = cmd.Password
		var userCredential credentials.User
		userCredential, responseError = credhubClient.SetUser(cmd.CredentialIdentifier, value, !cmd.NoOverwrite)
		output = interface{}(userCredential)

	} else if cmd.Type == "password" {
		var passwordCredential credentials.Password
		passwordCredential, responseError = credhubClient.SetPassword(cmd.CredentialIdentifier, values.Password(cmd.Password), !cmd.NoOverwrite)
		output = interface{}(passwordCredential)
	} else if cmd.Type == "json" {
		var jsonCredential credentials.JSON
		var unmarshalled values.JSON
		json.Unmarshal([]byte(cmd.Value), &unmarshalled)
		jsonCredential, responseError = credhubClient.SetJSON(cmd.CredentialIdentifier, unmarshalled, !cmd.NoOverwrite)
		output = interface{}(jsonCredential)
	} else {
		var valueCredential credentials.Value
		valueCredential, responseError = credhubClient.SetValue(cmd.CredentialIdentifier, values.Value(cmd.Value), !cmd.NoOverwrite)
		output = interface{}(valueCredential)
	}

	if responseError != nil {
		return nil, responseError
	}

	return &output, nil
}

func promptForInput(prompt string, value *string) {
	fmt.Printf(prompt)
	reader := bufio.NewReader(os.Stdin)
	val, _ := reader.ReadString('\n')
	*value = string(strings.TrimSpace(val))
}
