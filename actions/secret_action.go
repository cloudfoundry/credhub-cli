package actions

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type SecretAction struct {
	secretRepository repositories.SecretRepository
	config           config.Config
}

func NewSecretAction(secretRepository repositories.SecretRepository, config config.Config) SecretAction {
	return SecretAction{secretRepository: secretRepository, config: config}
}

func (action SecretAction) DoSecretAction(req *http.Request, secretIdentifier string) (models.Secret, error) {
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
