package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func (a *Api) Get(name string, id string) (models.CredentialResponse, error) {
	var (
		credential models.Printable
		err        error
	)

	repository := repositories.NewCredentialRepository(client.NewHttpClient(*a.Config))
	action := actions.NewAction(repository, a.Config)

	if name != "" {
		credential, err = action.DoAction(client.NewGetCredentialByNameRequest(*a.Config, name), name)
	} else if id != "" {
		credential, err = action.DoAction(client.NewGetCredentialByIdRequest(*a.Config, id), id)
	} else {
		return models.CredentialResponse{}, errors.NewMissingGetParametersError()
	}

	if err != nil {
		return models.CredentialResponse{}, err
	}

	return credential.(models.CredentialResponse), err
}

func (a *Api) GetByName(name string) (models.CredentialResponse, error) {
	return a.Get(name, "")
}

func (a *Api) GetById(id string) (models.CredentialResponse, error) {
	return a.Get("", id)
}
