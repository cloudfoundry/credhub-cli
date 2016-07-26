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
	response, err := DoSendRequest(r.httpClient, request)
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
