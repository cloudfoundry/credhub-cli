package actions

import (
	"encoding/json"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type Get struct {
	httpClient client.HttpClient
	config     config.Config
}

func NewGet(httpClient client.HttpClient, config config.Config) Get {
	return Get{httpClient: httpClient, config: config}
}

func (get Get) GetSecret(secretIdentifier string) (models.Secret, error) {
	err := config.ValidateConfig(get.config)

	if err != nil {
		return models.Secret{}, err
	}

	request := client.NewGetSecretRequest(get.config.ApiURL, secretIdentifier)

	response, err := get.httpClient.Do(request)
	if err != nil {
		return models.Secret{}, errors.NewNetworkError()
	}

	if response.StatusCode == 404 {
		return models.Secret{}, errors.NewSecretNotFoundError()
	}

	secretBody := new(models.SecretBody)

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(secretBody)
	if err != nil {
		return models.Secret{}, errors.NewResponseError()
	}

	return models.NewSecret(secretIdentifier, *secretBody), nil
}
