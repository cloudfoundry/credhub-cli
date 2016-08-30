package actions

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/errors"
	"github.com/pivotal-cf/credhub-cli/models"
)

func NewInfo(httpClient client.HttpClient, config config.Config) ServerInfo {
	return ServerInfo{httpClient: httpClient, config: config}
}

func (serverInfo ServerInfo) GetServerInfo() (models.Info, error) {
	request := client.NewInfoRequest(serverInfo.config)

	response, err := serverInfo.httpClient.Do(request)
	if err != nil {
		return models.Info{}, errors.NewNetworkError(err)
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
