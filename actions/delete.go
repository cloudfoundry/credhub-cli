package actions

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type Delete struct {
	secretRepository repositories.SecretRepository
	config           config.Config
}

func NewDelete(secretRepository repositories.SecretRepository, config config.Config) Delete {
	return Delete{secretRepository: secretRepository, config: config}
}

func (delete Delete) Delete(req *http.Request, secretIdentifier string) error {
	err := config.ValidateConfig(delete.config)

	if err != nil {
		return err
	}

	_, err = delete.secretRepository.SendRequest(req)

	return err
}
