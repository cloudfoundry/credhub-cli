package actions

import (
	"net/http"

	"reflect"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type Action struct {
	repository     repositories.Repository
	config         config.Config
	AuthRepository repositories.Repository
}

func NewAction(repository repositories.Repository, config config.Config) Action {
	action := Action{repository: repository, config: config}
	action.AuthRepository = repositories.NewAuthRepository(client.NewHttpClient(config.AuthURL))
	return action
}

func (action Action) DoAction(req *http.Request, identifier string) (models.Item, error) {
	err := config.ValidateConfig(action.config)

	if err != nil {
		return models.NewItem(), err
	}

	var secret models.Item
	secret, err = action.repository.SendRequest(req, identifier)
	if err != nil && reflect.DeepEqual(err, errors.NewUnauthorizedError()) {
		// refresh
		refresh_request := client.NewRefreshTokenRequest(action.config)
		refreshed_token, err := action.AuthRepository.SendRequest(refresh_request, "")

		if err != nil {
			return models.NewItem(), errors.NewRefreshError()
		}

		action.config.AccessToken = refreshed_token.(models.Token).AccessToken
		action.config.RefreshToken = refreshed_token.(models.Token).RefreshToken

		config.WriteConfig(action.config)

		req.Header.Set("Authorization", "Bearer "+action.config.AccessToken)
		secret, err = action.DoAction(req, identifier)
	} else if err != nil {
		return models.NewItem(), err
	}

	return secret, nil
}
