package actions

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories"
)

type Action struct {
	repository repositories.Repository
	config     config.Config
}

func NewAction(repository repositories.Repository, config config.Config) Action {
	return Action{repository: repository, config: config}
}

func (action Action) DoAction(req *http.Request, identifier string) (models.Item, error) {
	err := config.ValidateConfig(action.config)

	if err != nil {
		return models.NewItem(), err
	}

	secret, err := action.repository.SendRequest(req, identifier)
	if err != nil {
		return models.NewItem(), err
	}

	return secret, nil
}
