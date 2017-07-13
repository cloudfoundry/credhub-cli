package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func Delete(credentialIdentifier string) error {
	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, &cfg)

	_, err := action.DoAction(client.NewDeleteCredentialRequest(cfg, credentialIdentifier), credentialIdentifier)

	return err
}
