package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func (a *Api) Delete(credentialIdentifier string) error {
	repository := repositories.NewCredentialRepository(client.NewHttpClient(*a.Config))
	action := actions.NewAction(repository, a.Config)

	_, err := action.DoAction(client.NewDeleteCredentialRequest(*a.Config, credentialIdentifier), credentialIdentifier)

	return err
}
