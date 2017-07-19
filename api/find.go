package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func (a *Api) Find(partialCredentialIdentifier string, pathIdentifier string, allPaths bool) (models.CredentialResponse, error) {
	var credentials models.Printable
	var err error
	var repository repositories.Repository

	if allPaths {
		repository = repositories.NewAllPathRepository(client.NewHttpClient(*a.Config))
	} else {
		repository = repositories.NewCredentialQueryRepository(client.NewHttpClient(*a.Config))
	}

	action := actions.NewAction(repository, a.Config)

	if allPaths {
		credentials, err = action.DoAction(client.NewFindAllCredentialPathsRequest(*a.Config), "")
	} else if partialCredentialIdentifier != "" {
		credentials, err = action.DoAction(client.NewFindCredentialsBySubstringRequest(*a.Config, partialCredentialIdentifier), partialCredentialIdentifier)
	} else {
		credentials, err = action.DoAction(client.NewFindCredentialsByPathRequest(*a.Config, pathIdentifier), partialCredentialIdentifier)
	}

	if err != nil {
		return models.CredentialResponse{}, err
	}

	return credentials.(models.CredentialResponse), err
}
