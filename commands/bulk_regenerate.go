package commands

import (
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type BulkRegenerateCommand struct {
	SignedBy   string `required:"yes" long:"signed-by" description:"Selects the credential whose children should recursively be regenerated"`
	OutputJson bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
}

func (cmd BulkRegenerateCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	credhub, err := initializeCredhubClient(cfg)
	if err != nil {
		return err
	}

	credentials, err := credhub.BulkRegenerate(cmd.SignedBy)

	if err != nil {
		return err
	}

	printCredential(cmd.OutputJson, credentials)

	return nil
}
