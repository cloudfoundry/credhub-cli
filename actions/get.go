package actions

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type Get struct {
	secretRepository repositories.SecretRepository
	config           config.Config
}

func NewGet(secretRepository repositories.SecretRepository, config config.Config) Get {
	return Get{secretRepository: secretRepository, config: config}
}

func (get Get) GetSecret(req *http.Request, secretIdentifier string) (models.Secret, error) {
	err := config.ValidateConfig(get.config)

	if err != nil {
		return models.Secret{}, err
	}

	secretBody, err := get.secretRepository.SendRequest(req)
	if err != nil {
		return models.Secret{}, err
	}

	return models.NewSecret(secretIdentifier, secretBody), nil

}

//if response.StatusCode == 404 {
//return models.Secret{}, errors.ParseError(response.Body)
//}
