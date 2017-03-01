package repositories

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/errors"

	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"

	cm_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
)

type secretQueryRepository struct {
	httpClient client.HttpClient
}

func NewSecretQueryRepository(httpClient client.HttpClient) Repository {
	return secretQueryRepository{httpClient: httpClient}
}

func (r secretQueryRepository) SendRequest(request *http.Request, identifier string) (models.Printable, error) {
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return models.SecretQueryResponseBody{}, err
	}

	decoder := json.NewDecoder(response.Body)
	findResponseBody := models.SecretQueryResponseBody{}
	err = decoder.Decode(&findResponseBody)
	if err != nil {
		return models.SecretQueryResponseBody{}, cm_errors.NewResponseError()
	} else if len(findResponseBody.Credentials) < 1 {
		return models.SecretQueryResponseBody{}, errors.NewNoMatchingCredentialsFoundError()
	}
	return findResponseBody, nil
}
