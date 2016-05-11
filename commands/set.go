package commands

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type SetCommand struct {
	SecretIdentifier string `short:"n" long:"name" description:"Selects the secret being set"`
	SecretContent    string `short:"s" long:"secret" description:"Sets a value for a secret name"`
}

func (cmd SetCommand) Execute([]string) error {
	config := config.ReadConfig()

	request := client.NewPutSecretRequest(config.ApiURL, cmd.SecretIdentifier, cmd.SecretContent)

	response, _ := http.DefaultClient.Do(request)

	PrintResponse(response.Body)

	return nil
}
