package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type authRepository struct {
	httpClient client.HttpClient
}

func NewAuthRepository(httpClient client.HttpClient) Repository {
	return authRepository{httpClient: httpClient}
}

func (r authRepository) SendRequest(request *http.Request, identifier string) (models.Item, error) {
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return models.Token{}, err
	}

	decoder := json.NewDecoder(response.Body)
	token := models.Token{}
	err = decoder.Decode(&token)
	if err != nil {
		return models.Token{}, errors.NewResponseError()
	}
	return token, nil
}
