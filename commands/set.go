package commands

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
)

type SetCommand struct {
	SecretIdentifier string `short:"i" long:"identifier" description:"The Identifier of the Secret"`
	SecretContent    string `short:"s" long:"secret" description:"The Content of the Secret"`
}

func (cmd SetCommand) Execute([]string) error {
	request := client.NewPutSecretRequest(CM.ApiURL, cmd.SecretIdentifier, cmd.SecretContent)

	response, _ := http.DefaultClient.Do(request)

	PrintResponse(response.Body)

	return nil
}
