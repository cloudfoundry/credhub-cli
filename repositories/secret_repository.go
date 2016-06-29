package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type secretRepository struct {
	httpClient client.HttpClient
}

func NewSecretRepository(httpClient client.HttpClient) Repository {
	return secretRepository{httpClient: httpClient}
}

func (r secretRepository) SendRequest(request *http.Request, identifier string) (models.Item, error) {
	response, err := doSendRequest(r.httpClient, request)
	if err != nil {
		return models.Secret{}, err
	}

	if request.Method == "DELETE" {
		return models.Secret{}, nil
	}

	decoder := json.NewDecoder(response.Body)
	secretBody := models.SecretBody{}
	err = decoder.Decode(&secretBody)
	if err != nil {
		return models.Secret{}, errors.NewResponseError()
	}
	return models.NewSecret(identifier, secretBody), nil
}

func doSendRequest(client client.HttpClient, request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)

	if err != nil {
		return nil, errors.NewNetworkError()
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, errors.ParseError(response.Body)
	}
	return response, nil
}
