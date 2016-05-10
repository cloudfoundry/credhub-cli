package commands

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
)

type GetCommand struct {
	SecretIdentifier string `short:"n" long:"name" description:"Selects the secret to retrieve"`
}

func (cmd GetCommand) Execute([]string) error {
	request := client.NewGetSecretRequest(CM.ApiURL, cmd.SecretIdentifier)

	response, _ := http.DefaultClient.Do(request)

	PrintResponse(response.Body)

	return nil
}
