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
	secret := models.Secret{}
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return secret, err
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&secret.SecretBody)

	if err != nil {
		return secret, cm_errors.NewResponseError()
	}
	if len(secret.SecretBody["credentials"].([]interface{})) == 0 {
		return secret, errors.NewNoMatchingCredentialsFoundError()
	}

	return secret, nil
}
