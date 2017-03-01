package actions

import (
	"net/http"

	"reflect"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
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

func (action Action) DoAction(req *http.Request, identifier string) (models.Printable, error) {
	err := config.ValidateConfig(action.config)

	if err != nil {
		return nil, err
	}

	bodyClone := client.NewBodyClone(req)

	item, err := action.repository.SendRequest(req, identifier)

	if reflect.DeepEqual(err, errors.NewAccessTokenExpiredError()) {
		req.Body = bodyClone
		item, err = action.refreshTokenAndResendRequest(req, identifier)
	}

	return item, err
}

func (action Action) refreshTokenAndResendRequest(req *http.Request, identifier string) (models.Printable, error) {
	err := action.refreshToken()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+action.config.AccessToken)
	item, err := action.repository.SendRequest(req, identifier)
	if err != nil {
		return nil, err
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
