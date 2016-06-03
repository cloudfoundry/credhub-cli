package actions

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type Set struct {
	secretRepository repositories.SecretRepository
	config           config.Config
}

func NewSet(secretRepository repositories.SecretRepository, config config.Config) Set {
	return Set{
		secretRepository: secretRepository,
		config:           config,
	}
}

func (set Set) Set(req *http.Request, secretIdentifier string) (models.Secret, error) {
	err := config.ValidateConfig(set.config)

	if err != nil {
		return models.Secret{}, err
	}

	secretBody, err := set.secretRepository.SendRequest(req)
	if err != nil {
		return models.Secret{}, err
	}

	return models.NewSecret(secretIdentifier, secretBody), nil
}
