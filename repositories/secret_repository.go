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

type emptyBodyRepository struct {
	httpClient client.HttpClient
}

func NewSecretRepository(httpClient client.HttpClient) SecretRepository {
	return secretRepository{httpClient: httpClient}
}

func NewEmptyBodyRepository(httpClient client.HttpClient) SecretRepository {
	return emptyBodyRepository{httpClient: httpClient}
}

func (r secretRepository) SendRequest(request *http.Request) (models.SecretBody, error) {
	response, err := doSendRequest(r.httpClient, request)
	if err != nil {
		return models.SecretBody{}, err
	}

	decoder := json.NewDecoder(response.Body)
	secretBody := models.SecretBody{}
	err = decoder.Decode(&secretBody)
	if err != nil {
		return models.SecretBody{}, errors.NewResponseError()
	}
	return secretBody, nil
}

func (r emptyBodyRepository) SendRequest(request *http.Request) (models.SecretBody, error) {
	_, err := doSendRequest(r.httpClient, request)
	if err != nil {
		return models.SecretBody{}, err
	}

	return models.SecretBody{}, nil
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
