package commands

import (
	"fmt"

	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

type ImportCommand struct {
	File string `short:"f" long:"file" description:"File containing credentials to import. File must be in yaml format containing a list of credentials under the key 'credentials'. Name, type and value are required for each credential in the list." required:"true"`
}

func (cmd ImportCommand) Execute([]string) error {

	results, err := api.Import(cmd.File)
	if err != nil {
		return err
	}
	for _, result := range results {
		if result.Err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", result.Err)
		} else {
			models.Println(result.Cred, false)
		}
	}

	return nil
}
