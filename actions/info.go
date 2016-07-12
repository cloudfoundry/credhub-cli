package actions

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

type Version struct {
	httpClient client.HttpClient
	config     config.Config
}

func NewInfo(httpClient client.HttpClient, config config.Config) Version {
	return Version{httpClient: httpClient, config: config}
}

func (version Version) GetServerInfo() (models.Info, error) {
	request := client.NewInfoRequest(version.config.ApiURL)

	response, err := version.httpClient.Do(request)
	if err != nil {
		return models.Info{}, errors.NewNetworkError()
	}

	if response.StatusCode != http.StatusOK {
		return models.Info{}, errors.NewInvalidTargetError()
	}

	info := new(models.Info)

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(info)

	if err != nil {
		return models.Info{}, err
	}

	return *info, nil
}
