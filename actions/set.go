package actions

import (
	"github.com/pivotal-cf/cm-cli/client"
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

func (set Set) SetSecret(secretIdentifier string, value string, contentType string) (models.Secret, error) {
	err := config.ValidateConfig(set.config)

	if err != nil {
		return models.Secret{}, err
	}

	request := client.NewPutSecretRequest(set.config.ApiURL, secretIdentifier, value, contentType)

	secretBody, err := set.secretRepository.SendRequest(request)
	if err != nil {
		return models.Secret{}, err
	}

	return models.NewSecret(secretIdentifier, secretBody), nil
}
