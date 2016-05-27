package actions

import (
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/errors"
)

type Api struct {
	httpClient client.HttpClient
}

func NewApi(httpClient client.HttpClient) Api {
	return Api{httpClient: httpClient}
}

func (api Api) ValidateTarget(targetUrl string) error {
	request := client.NewInfoRequest(targetUrl)

	response, err := api.httpClient.Do(request)
	if err != nil {
		return errors.NewNetworkError()
	}

	if response.StatusCode != 200 {
		return errors.NewInvalidTargetError()
	}

	return nil
}
