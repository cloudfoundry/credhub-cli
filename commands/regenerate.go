package commands

import "github.com/cloudfoundry-incubator/credhub-cli/config"

type RegenerateCommand struct {
	CredentialIdentifier string `required:"yes" short:"n" long:"name" description:"Selects the credential to regenerate"`
	OutputJson           bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
}

func (cmd RegenerateCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	credhub, err := initializeCredhubClient(cfg)
	if err != nil {
		return err
	}

	credential, err := credhub.Regenerate(cmd.CredentialIdentifier)

	if err != nil {
		return err
	}

	printCredential(cmd.OutputJson, credential)

	return nil
}
