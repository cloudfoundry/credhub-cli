package commands

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
)

type SetCommand struct {
	SecretIdentifier string `short:"n" long:"name" description:"Selects the secret being set"`
	SecretContent    string `short:"s" long:"secret" description:"Sets a value for a secret name"`
}

func (cmd SetCommand) Execute([]string) error {
	request := client.NewPutSecretRequest(CM.ApiURL, cmd.SecretIdentifier, cmd.SecretContent)

	response, _ := http.DefaultClient.Do(request)

	PrintResponse(response.Body)

	return nil
}
