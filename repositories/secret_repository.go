package repositories

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/pivotal-cf/cm-cli/client"
	cli_errors "github.com/pivotal-cf/cm-cli/errors"
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
		return models.SecretBody{}, cli_errors.NewNetworkError()
	}

	decoder := json.NewDecoder(response.Body)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		serverError := models.ServerError{}
		decoder.Decode(&serverError)

		return models.SecretBody{}, errors.New(serverError.Message)
	}

	secretBody := models.SecretBody{}
	decoder.Decode(&secretBody)

	return secretBody, nil
}
