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

func (set Set) SetValue(secretIdentifier string, value string) (models.Secret, error) {
	err := config.ValidateConfig(set.config)

	if err != nil {
		return models.Secret{}, err
	}

	request := client.NewPutValueRequest(set.config.ApiURL, secretIdentifier, value)

	secretBody, err := set.secretRepository.SendRequest(request)
	if err != nil {
		return models.Secret{}, err
	}

	return models.NewSecret(secretIdentifier, secretBody), nil
}

func (set Set) SetCertificate(secretIdentifier string, ca string, pub string, priv string) (models.Secret, error) {
	err := config.ValidateConfig(set.config)

	if err != nil {
		return models.Secret{}, err
	}

	request := client.NewPutCertificateRequest(set.config.ApiURL, secretIdentifier, ca, pub, priv)

	secretBody, err := set.secretRepository.SendRequest(request)
	if err != nil {
		return models.Secret{}, err
	}

	return models.NewSecret(secretIdentifier, secretBody), nil
}
