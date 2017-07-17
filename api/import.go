package api

import (
	"net/http"
	"reflect"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func Import(file string) (results []struct {
	Name string
	Cred models.Printable
	Err  error
}, err error) {
	var name string
	var repository repositories.Repository
	var bulkImport models.CredentialBulkImport
	var request *http.Request

	err = bulkImport.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := config.ReadConfig()
	repository = repositories.NewCredentialRepository(client.NewHttpClient(cfg))
	action := actions.NewAction(repository, &cfg)

	for _, credential := range bulkImport.Credentials {
		var result struct {
			Name string
			Cred models.Printable
			Err  error
		}
		request = client.NewSetRequest(cfg, credential)

		switch credentialName := credential["name"].(type) {
		case string:
			name = credentialName
		default:
			name = ""
		}

		cred, err := action.DoAction(request, name)

		result.Name = name
		result.Cred = cred
		result.Err = err

		results = append(results, result)

		if err != nil {
			if isAuthenticationError(err) {
				return results, err
			}
		}

	}
	return results, nil
}

func isAuthenticationError(err error) bool {
	return reflect.DeepEqual(err, errors.NewNoApiUrlSetError()) ||
		reflect.DeepEqual(err, errors.NewRevokedTokenError()) ||
		reflect.DeepEqual(err, errors.NewRefreshError())
}
