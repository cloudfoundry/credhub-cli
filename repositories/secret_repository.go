package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	cm_errors "github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type secretRepository struct {
	httpClient client.HttpClient
}

func NewSecretRepository(httpClient client.HttpClient) Repository {
	return secretRepository{httpClient: httpClient}
}

func (r secretRepository) SendRequest(request *http.Request, identifier string) (models.Item, error) {
	response, err := DoSendRequest(r.httpClient, request)
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
		return models.Secret{}, cm_errors.NewResponseError()
	}
	return models.NewSecret(identifier, secretBody), nil
}
