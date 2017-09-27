package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
)

type GetCommand struct {
	Name             string `short:"n" long:"name" description:"Name of the credential to retrieve"`
	Id               string `long:"id" description:"ID of the credential to retrieve"`
	NumberOfVersions int    `long:"versions" description:"Number of versions of the credential to retrieve"`
	OutputJson       bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd GetCommand) Execute([]string) error {
	var (
		credential credentials.Credential
		err        error
	)

	cfg := config.ReadConfig()

	credhubClient, err := initializeCredhubClient(cfg)
	if err != nil {
		return err
	}

	var arrayOfCredentials []credentials.Credential

	if cmd.Name != "" {
		if cmd.NumberOfVersions != 0 {
			arrayOfCredentials, err = credhubClient.GetNVersions(cmd.Name, cmd.NumberOfVersions)
		} else {
			credential, err = credhubClient.GetLatestVersion(cmd.Name)
		}
	} else if cmd.Id != "" {
		credential, err = credhubClient.GetById(cmd.Id)
	} else {
		return errors.NewMissingGetParametersError()
	}

	if err != nil {
		return err
	}

	if arrayOfCredentials != nil {
		output := map[string][]credentials.Credential{
			"versions": arrayOfCredentials,
		}
		printCredential(cmd.OutputJson, output)
	} else {
		printCredential(cmd.OutputJson, credential)
	}

	return nil
}
