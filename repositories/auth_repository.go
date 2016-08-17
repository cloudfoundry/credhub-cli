package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type authRepository struct {
	httpClient         client.HttpClient
	expectResponseBody bool
}

func NewAuthRepository(httpClient client.HttpClient, expectResponseBody bool) Repository {
	return authRepository{httpClient: httpClient, expectResponseBody: expectResponseBody}
}

func (r authRepository) SendRequest(request *http.Request, identifier string) (models.Item, error) {
	response, err := DoSendRequest(r.httpClient, request)
	if err != nil {
		return models.Token{}, err
	}

	token := models.Token{}
	if r.expectResponseBody {
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&token)
		if err != nil {
			return models.Token{}, errors.NewResponseError()
		}
	}
	return token, nil
}
