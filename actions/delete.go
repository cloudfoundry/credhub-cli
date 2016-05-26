package actions

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
)

type Delete struct {
	httpClient client.HttpClient
	config     config.Config
}

func NewDelete(httpClient client.HttpClient, config config.Config) Delete {
	return Delete{httpClient: httpClient, config: config}
}

func (delete Delete) Delete(secretIdentifier string) error {
	err := config.ValidateConfig(delete.config)

	if err != nil {
		return err
	}

	request := client.NewDeleteSecretRequest(delete.config.ApiURL, secretIdentifier)

	response, err := delete.httpClient.Do(request)
	if err != nil {
		return errors.NewNetworkError()
	}

	if response.StatusCode == http.StatusNotFound {
		return errors.NewSecretNotFoundError()
	} else if response.StatusCode != http.StatusOK {
		return errors.NewSecretBadRequestError()
	}

	return nil
}
