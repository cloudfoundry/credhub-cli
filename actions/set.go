package actions

import (
	"encoding/json"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
)

type Set struct {
	httpClient HttpClient
	config     config.Config
}

func NewSet(httpClient HttpClient, config config.Config) Set {
	return Set{httpClient: httpClient, config: config}
}

func (set Set) SetSecret(secretIdentifier string, value string) (client.Secret, error) {
	err := config.ValidateConfig(set.config)

	if err != nil {
		return client.Secret{}, err
	}

	request := client.NewPutSecretRequest(set.config.ApiURL, secretIdentifier, value)

	response, err := set.httpClient.Do(request)

	if err != nil {
		return client.Secret{}, errors.NewNetworkError()
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return client.Secret{}, errors.NewInvalidStatusError()
	}

	secretBody := new(client.SecretBody)

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(secretBody)
	if err != nil {
		return client.Secret{}, errors.NewResponseError()
	}

	secret := client.NewSecret(secretIdentifier, *secretBody)

	return secret, nil
}
