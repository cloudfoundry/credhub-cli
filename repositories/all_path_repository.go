package repositories

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/errors"

	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"

	cm_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
)

type allPathRepository struct {
	httpClient client.HttpClient
}

func NewAllPathRepository(httpClient client.HttpClient) Repository {
	return allPathRepository{httpClient: httpClient}
}

func (r allPathRepository) SendRequest(request *http.Request, ignoredIdentifier string) (models.Printable, error) {
	secret := models.CredentialResponse{}

	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return secret, err
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&secret.ResponseBody)

	if err != nil {
		return secret, cm_errors.NewResponseError()
	}

	if len(secret.ResponseBody["paths"].([]interface{})) == 0 {
		return secret, errors.NewNoMatchingCredentialsFoundError()
	}
	return secret, nil
}
