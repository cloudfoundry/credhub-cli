package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func (a *Api) Regenerate(credentialIdentifier string) (models.CredentialResponse, error) {
	repository := repositories.NewCredentialRepository(client.NewHttpClient(*a.Config))
	action := actions.NewAction(repository, a.Config)

	credential, err := action.DoAction(client.NewRegenerateCredentialRequest(*a.Config, credentialIdentifier), credentialIdentifier)
	if err != nil {
		return models.CredentialResponse{}, err
	}

	return credential.(models.CredentialResponse), err
}
