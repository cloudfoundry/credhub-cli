package commands

import (
	"net/http"

	"fmt"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type DeleteCommand struct {
	SecretIdentifier string `short:"n" long:"name" description:"Selects the secret to delete"`
}

func (cmd DeleteCommand) Execute([]string) error {
	config := config.ReadConfig()

	request := client.NewDeleteSecretRequest(config.ApiURL, cmd.SecretIdentifier)

	http.DefaultClient.Do(request)

	fmt.Println("Secret successfully deleted")

	return nil
}
