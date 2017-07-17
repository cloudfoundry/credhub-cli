package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func Find(partialCredentialIdentifier string, pathIdentifier string, allPaths bool) (models.CredentialResponse, error) {
	var credentials models.Printable
	var err error
	var repository repositories.Repository

	cfg := config.ReadConfig()

	if allPaths {
		repository = repositories.NewAllPathRepository(client.NewHttpClient(cfg))
	} else {
		repository = repositories.NewCredentialQueryRepository(client.NewHttpClient(cfg))
	}

	action := actions.NewAction(repository, &cfg)

	if allPaths {
		credentials, err = action.DoAction(client.NewFindAllCredentialPathsRequest(cfg), "")
	} else if partialCredentialIdentifier != "" {
		credentials, err = action.DoAction(client.NewFindCredentialsBySubstringRequest(cfg, partialCredentialIdentifier), partialCredentialIdentifier)
	} else {
		credentials, err = action.DoAction(client.NewFindCredentialsByPathRequest(cfg, pathIdentifier), partialCredentialIdentifier)
	}

	return credentials.(models.CredentialResponse), err
}
