package commands

import (
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type GetCommand struct {
	SecretIdentifier string `short:"n" long:"name" description:"Selects the secret to retrieve"`
}

func (cmd GetCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	err := config.ValidateConfig(cfg)
	if err != nil {
		return err
	}

	request := client.NewGetSecretRequest(cfg.ApiURL, cmd.SecretIdentifier)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return NewNetworkError()
	}

	if response.StatusCode == 404 {
		return NewSecretNotFoundError()
	}

	responseMsg, _ := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	return NewResponseError()
	//}

	secretBody := new(client.SecretBody)

	if json.Unmarshal(responseMsg, &secretBody) != nil {
		return NewResponseError()
	}

	secret := client.NewSecret(cmd.SecretIdentifier, *secretBody)

	secret.PrintSecret()

	return nil
}
