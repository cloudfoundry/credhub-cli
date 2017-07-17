package api

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func Import(file string) (results []struct {
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

		result.Cred = cred
		result.Err = err

		results = append(results, result)
	}
	return results, nil
}
