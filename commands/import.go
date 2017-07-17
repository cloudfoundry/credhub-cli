package commands

import (
	"fmt"

	"net/http"

	"os"

	"reflect"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
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
	var name string

	for i, credential := range bulkImport.Credentials {

		switch credentialName := credential["name"].(type) {
		case string:
			name = credentialName
		default:
			name = ""
		}

		credentialType := credential["type"].(string)
		value := credential["value"]
		overwrite := credential["overwrite"].(bool)

		result, err := api.Set(credentialType, name, value, overwrite)

		if err != nil {
			if isAuthenticationError(err) {
				fmt.Println(err)
				break
			}
			fmt.Fprintf(os.Stderr, "Credential '%s' at index %d could not be set: %v\n", name, i, err)
			continue
		}

		models.Println(result, false)
	}
}
func isAuthenticationError(err error) bool {
	return reflect.DeepEqual(err, errors.NewNoApiUrlSetError()) ||
		reflect.DeepEqual(err, errors.NewRevokedTokenError()) ||
		reflect.DeepEqual(err, errors.NewRefreshError())
}
