package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func Set(credentialType string, name string, value interface{}, overwrite bool) (models.CredentialResponse, error) {
	cfg := config.ReadConfig()

	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))

	action := actions.NewAction(repository, &cfg)

	request := client.NewSetCredentialRequest(cfg, credentialType, name, value, overwrite)

	result, err := action.DoAction(request, name)

	if err != nil {
		return models.CredentialResponse{}, err
	}

	return result.(models.CredentialResponse), err
}
