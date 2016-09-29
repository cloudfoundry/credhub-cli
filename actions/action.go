package actions

import (
	"net/http"

	"reflect"

	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/errors"
	"github.com/pivotal-cf/credhub-cli/models"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type Action struct {
	repository     repositories.Repository
	config         config.Config
	AuthRepository repositories.Repository
}

func NewAction(repository repositories.Repository, config config.Config) Action {
	action := Action{repository: repository, config: config}
	action.AuthRepository = repositories.NewAuthRepository(client.NewHttpClient(config), true)
	return action
}

func (action Action) DoAction(req *http.Request, identifier string) (interface{}, error) {
	err := config.ValidateConfig(action.config)

	if err != nil {
		return struct {}{}, err
	}

	bodyClone := client.NewBodyClone(req)

	item, err := action.repository.SendRequest(req, identifier)

	if reflect.DeepEqual(err, errors.NewUnauthorizedError()) {
		req.Body = bodyClone
		item, err = action.refreshTokenAndResendRequest(req, identifier)
	}

	if err != nil {
		return struct {}{}, err
	}

	return item, nil
}

func (action Action) refreshTokenAndResendRequest(req *http.Request, identifier string) (interface{}, error) {
	err := action.refreshToken()
	if err != nil {
		return struct {}{}, err
	}

	req.Header.Set("Authorization", "Bearer "+action.config.AccessToken)
	item, err := action.repository.SendRequest(req, identifier)
	if err != nil {
		return struct {}{}, err
	}

	return item, nil
}

func (action *Action) refreshToken() error {
	refresh_request := client.NewRefreshTokenRequest(action.config)
	refreshed_token, err := action.AuthRepository.SendRequest(refresh_request, "")

	if err != nil {
		return errors.NewRefreshError()
	}

	action.config.AccessToken = refreshed_token.(models.Token).AccessToken
	action.config.RefreshToken = refreshed_token.(models.Token).RefreshToken

	config.WriteConfig(action.config)

	return nil
}
