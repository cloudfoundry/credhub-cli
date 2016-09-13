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
	SecretIdentifier string   `short:"n" required:"yes" long:"name" description:"Selects the secret being set"`
	ContentType      string   `short:"t" long:"type" description:"Sets the type of secret to store or generate. Default: 'value'"`
	NoOverwrite      bool     `short:"O" long:"no-overwrite" description:"Credential is not modified if stored value already exists"`
	Length           int      `short:"l" long:"length" description:"Sets length of generated value (Default: 20)"`
	ExcludeSpecial   bool     `short:"S" long:"exclude-special" description:"Exclude special characters from generated value"`
	ExcludeNumber    bool     `short:"N" long:"exclude-number" description:"Exclude number characters from generated value"`
	ExcludeUpper     bool     `short:"U" long:"exclude-upper" description:"Exclude upper alpha characters from generated value"`
	ExcludeLower     bool     `short:"L" long:"exclude-lower" description:"Exclude lower alpha characters from generated value"`
	CommonName       string   `short:"c" long:"common-name" description:"Sets the common name of the generated certificate"`
	Organization     string   `short:"o" long:"organization" description:"Sets the organization of the generated certificate"`
	OrganizationUnit string   `short:"u" long:"organization-unit" description:"Sets the organization unit of the generated certificate"`
	Locality         string   `short:"i" long:"locality" description:"Sets the locality/city of the generated certificate"`
	State            string   `short:"s" long:"state" description:"Sets the state/province of the generated certificate"`
	Country          string   `short:"y" long:"country" description:"Sets the country of the generated certificate"`
	AlternativeName  []string `short:"a" long:"alternative-name" description:"Sets an alternative name of the generated certificate. Multiple alternative names can be set"`
	KeyLength        int      `short:"k" long:"key-length" description:"Sets the bit length of the generated key"`
	Duration         int      `short:"d" long:"duration" description:"Sets the valid duration (in days) for the generated certificate"`
	Ca               string   `long:"ca" description:"Selects the CA used to sign the generated certificate"`
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
