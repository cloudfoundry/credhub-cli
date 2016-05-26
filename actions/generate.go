package actions

import (
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type Generate struct {
	secretRepository repositories.SecretRepository
	config           config.Config
}

func NewGenerate(secretRepository repositories.SecretRepository, config config.Config) Generate {
	return Generate{
		secretRepository: secretRepository,
		config:           config,
	}
}

func (g Generate) GenerateSecret(secretIdentifier string, parameters models.SecretParameters) (models.Secret, error) {
	err := config.ValidateConfig(g.config)

	if err != nil {
		return models.Secret{}, err
	}

	request := client.NewGenerateSecretRequest(g.config.ApiURL, secretIdentifier, parameters)

	secretBody, err := g.secretRepository.SendRequest(request)
	if err != nil {
		return models.Secret{}, err
	}

	return models.NewSecret(secretIdentifier, secretBody), nil
}
