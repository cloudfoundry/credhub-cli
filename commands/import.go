package commands

import (
	"fmt"

	"net/http"

	"os"

	"reflect"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

type ImportCommand struct {
	File string `short:"f" long:"file" description:"File containing credentials to import. File must be in yaml format containing a list of credentials under the key 'credentials'. Name, type and value are required for each credential in the list." required:"true"`
}

var (
	err        error
	repository repositories.Repository
	bulkImport models.CredentialBulkImport
	request    *http.Request
)

func (cmd ImportCommand) Execute([]string) error {
	err = bulkImport.ReadFile(cmd.File)

	if err != nil {
		return err
	}

	setCredentials(bulkImport)

	return nil
}

func setCredentials(bulkImport models.CredentialBulkImport) {
	var (
		name       string
		successful int
		failed     int
	)
	errors := make([]string, 0)

	cfg := config.ReadConfig()
	repository = repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, &cfg)

	for i, credential := range bulkImport.Credentials {
		request = client.NewSetRequest(cfg, credential)

		switch credentialName := credential["name"].(type) {
		case string:
			name = credentialName
		default:
			name = ""
		}

		result, err := action.DoAction(request, name)

		if err != nil {
			if isAuthenticationError(err) {
				fmt.Println(err)
				return
			}
			failure := fmt.Sprintf("Credential '%s' at index %d could not be set: %v", name, i, err)
			fmt.Println(failure + "\n")
			errors = append(errors, " - "+failure)
			failed++
			continue
		} else {
			successful++
		}
		models.Println(result, false)
	}

	fmt.Println("Import complete.")
	fmt.Fprintf(os.Stdout, "Successfully set: %d\n", successful)
	fmt.Fprintf(os.Stdout, "Failed to set: %d\n", failed)
	for _, v := range errors {
		fmt.Println(v)
	}
}
func isAuthenticationError(err error) bool {
	return reflect.DeepEqual(err, errors.NewNoApiUrlSetError()) ||
		reflect.DeepEqual(err, errors.NewRevokedTokenError()) ||
		reflect.DeepEqual(err, errors.NewRefreshError())
}
