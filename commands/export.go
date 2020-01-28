package commands

import (
	"fmt"
	"io/ioutil"

	"code.cloudfoundry.org/credhub-cli/config"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/models"
)

type ExportCommand struct {
	Path       string `short:"p" long:"path" description:"Path of credentials to export" required:"false"`
	File       string `short:"f" long:"file" description:"File in which to write credentials" required:"false"`
	OutputJSON bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
}

func (cmd ExportCommand) Execute([]string) error {
	allCredentials, err := getAllCredentialsForPath(cmd.Path)

	if err != nil {
		return err
	}

	exportCreds, err := models.ExportCredentials(allCredentials, cmd.OutputJSON)

	if err != nil {
		return err
	}

	if cmd.File == "" {
		fmt.Printf("%s", exportCreds)

		return err
	} else {
		return ioutil.WriteFile(cmd.File, exportCreds.Bytes, 0644)
	}
}

func getAllCredentialsForPath(path string) ([]credentials.Credential, error) {
	cfg := config.ReadConfig()
	credhubClient, err := initializeCredhubClient(cfg)

	if err != nil {
		return nil, err
	}

	allPaths, err := credhubClient.FindByPath(path)

	if err != nil {
		return nil, err
	}

	credentials := make([]credentials.Credential, len(allPaths.Credentials))
	for i, baseCred := range allPaths.Credentials {
		credential, err := credhubClient.GetLatestVersion(baseCred.Name)

		if err != nil {
			return nil, err
		}

		if credential.Type == "certificate" {
			certMetadata, err := credhubClient.GetCertificateMetadataByName(credential.Name)

			if err != nil {
				return nil, err
			}
			signedBy := certMetadata.SignedBy

			if signedBy != "" && signedBy != credential.Name {
				if cert, ok := credential.Value.(map[string]interface{}); ok {
					cert["ca"] = signedBy
					credential.Value = cert
				}
			}
		}
		credentials[i] = credential
	}

	return credentials, nil
}
