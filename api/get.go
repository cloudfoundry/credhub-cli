package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func Get(name string, id string) (models.Printable, error) {
	var (
		credential models.Printable
		err        error
	)

	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, &cfg)

	if name != "" {
		credential, err = action.DoAction(client.NewGetCredentialByNameRequest(cfg, name), name)
	} else if id != "" {
		credential, err = action.DoAction(client.NewGetCredentialByIdRequest(cfg, id), id)
	} else {
		return credential, errors.NewMissingGetParametersError()
	}

	if err != nil {
		return credential, err
	}

	return credential, err
}
