package actions

import (
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
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

func (set Set) GenerateSecret(secretIdentifier string) (client.Secret, error) {
	err := config.ValidateConfig(set.config)

	if err != nil {
		return client.Secret{}, err
	}

	request := client.NewGenerateSecretRequest(set.config.ApiURL, secretIdentifier)

	secretBody, err := set.secretRepository.SendRequest(request)
	if err != nil {
		return client.Secret{}, err
	}

	return client.NewSecret(secretIdentifier, secretBody), nil
}

func (set Set) SetSecret(secretIdentifier string, value string) (client.Secret, error) {
	err := config.ValidateConfig(set.config)

	if err != nil {
		return client.Secret{}, err
	}

	request := client.NewPutSecretRequest(set.config.ApiURL, secretIdentifier, value)

	secretBody, err := set.secretRepository.SendRequest(request)
	if err != nil {
		return client.Secret{}, err
	}

	return client.NewSecret(secretIdentifier, secretBody), nil
}
