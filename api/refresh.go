package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func (a *Api) Refresh() (models.Token, error) {
	request := client.NewRefreshTokenRequest(*a.Config)
	repository := repositories.NewAuthRepository(client.NewHttpClient(*a.Config), true)
	token, err := repository.SendRequest(request, "")

	if err != nil {
		return models.Token{}, err
	}

	a.Config.AccessToken = token.(models.Token).AccessToken
	a.Config.RefreshToken = token.(models.Token).RefreshToken

	return token.(models.Token), nil
}
