package actions

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type Action struct {
	secretRepository repositories.SecretRepository
	config           config.Config
}

func NewAction(secretRepository repositories.SecretRepository, config config.Config) Action {
	return Action{secretRepository: secretRepository, config: config}
}

func (action Action) DoAction(req *http.Request, secretIdentifier string) (models.Secret, error) {
	err := config.ValidateConfig(action.config)

	if err != nil {
		return models.Secret{}, err
	}

	secretBody, err := action.secretRepository.SendRequest(req)
	if err != nil {
		return models.Secret{}, err
	}

	return models.NewSecret(secretIdentifier, secretBody), nil
}
