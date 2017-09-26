package commands

import (
	"fmt"

	"net/http"

	"os"

	"reflect"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type ImportCommand struct {
	File string `short:"f" long:"file" description:"File containing credentials to import" required:"true"`
}

var (
	err        error
	bulkImport models.CredentialBulkImport
	request    *http.Request
)

func (cmd ImportCommand) Execute([]string) error {
	err = bulkImport.ReadFile(cmd.File)

	if err != nil {
		return err
	}

	err := setCredentials(bulkImport)

	return err
}

func setCredentials(bulkImport models.CredentialBulkImport) error {
	var (
		name       string
		successful int
		failed     int
	)
	errors := make([]string, 0)

	cfg := config.ReadConfig()

	var credhubClient *credhub.CredHub

	if clientCredentialsInEnvironment() {
		credhubClient, err = newCredhubClient(&cfg, os.Getenv("CREDHUB_CLIENT"), os.Getenv("CREDHUB_SECRET"), true)
	} else {
		credhubClient, err = newCredhubClient(&cfg, config.AuthClient, config.AuthPassword, false)
	}
	if err != nil {
		return err
	}

	err = config.ValidateConfig(cfg)
	if err != nil {
		if !clientCredentialsInEnvironment() || config.ValidateConfigApi(cfg) != nil {
			return err
		}
	}

	for i, credential := range bulkImport.Credentials {
		switch credentialName := credential["name"].(type) {
		case string:
			name = credentialName
		default:
			name = ""
		}

		result, err := credhubClient.SetCredential(name, credential["type"].(string), credential["value"], true)

		if err != nil {
			if isAuthenticationError(err) {
				return err
			}
			failure := fmt.Sprintf("Credential '%s' at index %d could not be set: %v", name, i, err)
			fmt.Println(failure + "\n")
			errors = append(errors, " - "+failure)
			failed++
			continue
		} else {
			successful++
		}
		printCredential(false, result)
	}

	fmt.Println("Import complete.")
	fmt.Fprintf(os.Stdout, "Successfully set: %d\n", successful)
	fmt.Fprintf(os.Stdout, "Failed to set: %d\n", failed)
	for _, v := range errors {
		fmt.Println(v)
	}

	return nil
}

func isAuthenticationError(err error) bool {
	return reflect.DeepEqual(err, errors.NewNoApiUrlSetError()) ||
		reflect.DeepEqual(err, errors.NewRevokedTokenError()) ||
		reflect.DeepEqual(err, errors.NewRefreshError())
}
