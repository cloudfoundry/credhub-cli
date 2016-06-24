package actions

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type CaAction struct {
	caRepository repositories.CaRepository
	config       config.Config
}

func NewCaAction(caRepository repositories.CaRepository, config config.Config) CaAction {
	return CaAction{caRepository: caRepository, config: config}
}

func (action CaAction) DoCaAction(req *http.Request, caIdentifier string) (models.Ca, error) {
	err := config.ValidateConfig(action.config)

	if err != nil {
		return models.Ca{}, err
	}

	caBody, err := action.caRepository.SendRequest(req)
	if err != nil {
		return models.Ca{}, err
	}

	return models.NewCa(caIdentifier, caBody), nil
}
