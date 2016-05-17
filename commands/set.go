package commands

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"io/ioutil"
)

type SetCommand struct {
	SecretIdentifier string `short:"n" long:"name" description:"Selects the secret being set"`
	SecretContent    string `short:"s" long:"secret" description:"Sets a value for a secret name"`
}

func (cmd SetCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	err := config.ValidateConfig(cfg)
	if err != nil {
		return err
	}

	request := client.NewPutSecretRequest(cfg.ApiURL, cmd.SecretIdentifier, cmd.SecretContent)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return NewNetworkError()
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return NewInvalidStatusError()
	}

	responseMsg, _ :=ioutil.ReadAll(response.Body)
	//if err != nil {
		//return NewResponseError()
	//}

	PrintResponse(responseMsg)

	return nil
}
