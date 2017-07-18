package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func (a *Api) Set(credentialType string, name string, value interface{}, overwrite bool) (models.CredentialResponse, error) {
	repository := repositories.NewCredentialRepository(client.NewHttpClient(*a.Config))

	action := actions.NewAction(repository, a.Config)

	request := client.NewSetCredentialRequest(*a.Config, credentialType, name, value, overwrite)

	result, err := action.DoAction(request, name)

	if err != nil {
		return models.CredentialResponse{}, err
	}

	return result.(models.CredentialResponse), err
}
