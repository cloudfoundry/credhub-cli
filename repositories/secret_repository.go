package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type SecretRepository interface {
	SendRequest(request *http.Request) (models.SecretBody, error)
}

type secretRepository struct {
	httpClient client.HttpClient
}

func NewSecretRepository(httpClient client.HttpClient) SecretRepository {
	return secretRepository{httpClient: httpClient}
}

func (r secretRepository) SendRequest(request *http.Request) (models.SecretBody, error) {
	response, err := r.httpClient.Do(request)

	if err != nil {
		return models.SecretBody{}, errors.NewNetworkError()
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return models.SecretBody{}, errors.NewInvalidStatusError()
	}

	secretBody := models.SecretBody{}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&secretBody)
	if err != nil {
		return models.SecretBody{}, errors.NewResponseError()
	}

	return secretBody, nil
}
