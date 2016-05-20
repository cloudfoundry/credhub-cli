package actions

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type Version struct {
	httpClient HttpClient
	config     config.Config
}

func NewVersion(httpClient HttpClient, config config.Config) Version {
	return Version{httpClient: httpClient, config: config}
}

func (version Version) GetServerVersion() string {
	cmVersion := "Not Found"
	request := client.NewInfoRequest(version.config.ApiURL)

	response, err := version.httpClient.Do(request)
	if err == nil && response.StatusCode == http.StatusOK {
		info := new(client.Info)

		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(info)

		if err == nil {
			cmVersion = info.App.Version
		}
	}

	return cmVersion
}
