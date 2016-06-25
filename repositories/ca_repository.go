package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type caRepository struct {
	httpClient client.HttpClient
}

func NewCaRepository(httpClient client.HttpClient) Repository {
	return caRepository{httpClient: httpClient}
}

func (r caRepository) SendRequest(request *http.Request, caIdentifier string) (models.Item, error) {
	response, err := r.doSendRequest(r.httpClient, request)
	if err != nil {
		return models.Ca{}, err
	}

	decoder := json.NewDecoder(response.Body)
	caBody := models.CaBody{}
	err = decoder.Decode(&caBody)
	if err != nil {
		return models.Ca{}, errors.NewResponseError()
	}
	return models.NewCa(caIdentifier, caBody), nil
}

func (r caRepository) doSendRequest(client client.HttpClient, request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)

	if err != nil {
		return nil, errors.NewNetworkError()
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, errors.ParseError(response.Body)
	}
	return response, nil
}
