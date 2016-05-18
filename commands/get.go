package commands

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"io/ioutil"
	"encoding/json"
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

	secret := new(client.Secret)

	if json.Unmarshal(responseMsg, &secret) != nil {
		return NewResponseError()
	}

	secret.PrintSecret()

	return nil
}
