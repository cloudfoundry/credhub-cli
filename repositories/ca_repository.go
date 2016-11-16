package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/errors"
	"github.com/pivotal-cf/credhub-cli/models"
)

type caRepository struct {
	httpClient client.HttpClient
}

func NewCaRepository(httpClient client.HttpClient) Repository {
	return caRepository{httpClient: httpClient}
}

func (r caRepository) SendRequest(request *http.Request, caIdentifier string) (models.Printable, error) {
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return models.Ca{}, err
	}

	decoder := json.NewDecoder(response.Body)
	decoded := map[string]interface{}{}

	err = decoder.Decode(&decoded)

	if err != nil {
		return models.Ca{}, errors.NewResponseError()
	}

	caBody := models.NewCaBody(decoded)

	return models.NewCa(caIdentifier, caBody), nil
}
