package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type CaRepository interface {
	SendRequest(request *http.Request) (models.CaBody, error)
}

type caRepository struct {
	httpClient client.HttpClient
}

func NewCaRepository(httpClient client.HttpClient) CaRepository {
	return caRepository{httpClient: httpClient}
}

func (r caRepository) SendRequest(request *http.Request) (models.CaBody, error) {
	response, err := r.doSendRequest(r.httpClient, request)
	if err != nil {
		return models.CaBody{}, err
	}

	decoder := json.NewDecoder(response.Body)
	caBody := models.CaBody{}
	err = decoder.Decode(&caBody)
	if err != nil {
		return models.CaBody{}, errors.NewResponseError()
	}
	return caBody, nil
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
