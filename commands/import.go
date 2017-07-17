package commands

import (
	"fmt"
	"reflect"

	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type ImportCommand struct {
	File string `short:"f" long:"file" description:"File containing credentials to import. File must be in yaml format containing a list of credentials under the key 'credentials'. Name, type and value are required for each credential in the list." required:"true"`
}

func (cmd ImportCommand) Execute([]string) error {

	results, err := api.Import(cmd.File)
	if err != nil {
		if isAuthenticationError(err) {
			fmt.Println(err)
		}
		return err
	}

	for i, result := range results {
		if result.Err != nil {
			fmt.Fprintf(os.Stderr, "Credential '%s' at index %d could not be set: %v\n", result.Name, i, result.Err)
		} else {
			models.Println(result.Cred, false)
		}
	}

	return nil
}

func isAuthenticationError(err error) bool {
	return reflect.DeepEqual(err, errors.NewNoApiUrlSetError()) ||
		reflect.DeepEqual(err, errors.NewRevokedTokenError()) ||
		reflect.DeepEqual(err, errors.NewRefreshError())
}
