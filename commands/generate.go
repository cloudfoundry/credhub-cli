package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/models"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type GenerateCommand struct {
	SecretIdentifier string   `short:"n" required:"yes" long:"name" description:"Name of the credential to generate"`
	ContentType      string   `short:"t" long:"type" description:"Sets the credential type to generate (Default: 'password')"`
	Length           int      `short:"l" long:"length" description:"[Password] Length of the generated value (Default: 20)"`
	ExcludeSpecial   bool     `short:"S" long:"exclude-special" description:"[Password] Exclude special characters from the generated value"`
	ExcludeNumber    bool     `short:"N" long:"exclude-number" description:"[Password] Exclude number characters from the generated value"`
	ExcludeUpper     bool     `short:"U" long:"exclude-upper" description:"[Password] Exclude upper alpha characters from the generated value"`
	ExcludeLower     bool     `short:"L" long:"exclude-lower" description:"[Password] Exclude lower alpha characters from the generated value"`
	Ca               string   `long:"ca" description:"[Certificate] Name of CA used to sign the generated certificate (Default: 'default')"`
	Duration         int      `short:"d" long:"duration" description:"[Certificate] Valid duration (in days) of the generated certificate (Default: 365)"`
	KeyLength        int      `short:"k" long:"key-length" description:"[Certificate, SSH, RSA] Bit length of the generated key (Default: 2048)"`
	CommonName       string   `short:"c" long:"common-name" description:"[Certificate] Common name of the generated certificate"`
	Organization     string   `short:"o" long:"organization" description:"[Certificate] Organization of the generated certificate"`
	OrganizationUnit string   `short:"u" long:"organization-unit" description:"[Certificate] Organization unit of the generated certificate"`
	Locality         string   `short:"i" long:"locality" description:"[Certificate] Locality/city of the generated certificate"`
	State            string   `short:"s" long:"state" description:"[Certificate] State/province of the generated certificate"`
	Country          string   `short:"y" long:"country" description:"[Certificate] Country of the generated certificate"`
	AlternativeName  []string `short:"a" long:"alternative-name" description:"[Certificate] A subject alternative name of the generated certificate (may be specified multiple times)"`
	NoOverwrite      bool     `short:"O" long:"no-overwrite" description:"Credential is not modified if stored value already exists"`
}

func (cmd GenerateCommand) Execute([]string) error {
	if cmd.ContentType == "" {
		cmd.ContentType = "password"
	}

	cfg := config.ReadConfig()
	repository := repositories.NewSecretRepository(client.NewHttpClient(cfg))

	parameters := models.SecretParameters{
		ExcludeSpecial:   cmd.ExcludeSpecial,
		ExcludeNumber:    cmd.ExcludeNumber,
		ExcludeUpper:     cmd.ExcludeUpper,
		ExcludeLower:     cmd.ExcludeLower,
		Length:           cmd.Length,
		CommonName:       cmd.CommonName,
		Organization:     cmd.Organization,
		OrganizationUnit: cmd.OrganizationUnit,
		Locality:         cmd.Locality,
		State:            cmd.State,
		Country:          cmd.Country,
		AlternativeName:  cmd.AlternativeName,
		KeyLength:        cmd.KeyLength,
		Duration:         cmd.Duration,
		Ca:               cmd.Ca,
	}

	action := actions.NewAction(repository, cfg)
	request := client.NewGenerateSecretRequest(cfg, cmd.SecretIdentifier, parameters, cmd.ContentType, !cmd.NoOverwrite)
	secret, err := action.DoAction(request, cmd.SecretIdentifier)

	if err != nil {
		return err
	}

	fmt.Println(secret)

	return nil
}
