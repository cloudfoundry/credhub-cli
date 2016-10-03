package repositories

import (
	"net/http"

	"github.com/pivotal-cf/credhub-cli/errors"

	"encoding/json"

	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/models"

	cm_errors "github.com/pivotal-cf/credhub-cli/errors"
)

type allPathRepository struct {
	httpClient client.HttpClient
}

func NewAllPathRepository(httpClient client.HttpClient) Repository {
	return allPathRepository{httpClient: httpClient}
}

func (r allPathRepository) SendRequest(request *http.Request, ignoredIdentifier string) (interface{}, error) {
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return models.AllPathResponseBody{}, err
	}

	decoder := json.NewDecoder(response.Body)
	findResponseBody := models.AllPathResponseBody{}
	err = decoder.Decode(&findResponseBody)
	if err != nil {
		return models.AllPathResponseBody{}, cm_errors.NewResponseError()
	} else if len(findResponseBody.Paths) < 1 {
		return models.AllPathResponseBody{}, errors.NewNoMatchingCredentialsFoundError()
	}
	return findResponseBody, nil
}
