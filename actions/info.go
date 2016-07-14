package actions

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
	"fmt"
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

func NewToken(httpClient client.HttpClient, config config.Config) Version {
	return Version{httpClient: httpClient, config: config}
}

func (version Version) GetToken(user string, pass string) (models.Token, error) {
	request := client.NewTokenRequest(version.config.AuthURL, user, pass)
	//fmt.Println("Auth URL is #%v", version.config.AuthURL)
	//body, _ := ioutil.ReadAll(request.Body)
	fmt.Println("starting request with config: %#v", request)
	response, err := version.httpClient.Do(request)
	fmt.Println("finished request with response: %#v", response)
	if err != nil {
		return models.Token{}, errors.NewNetworkError()
	}

	if response.StatusCode != http.StatusOK {
		return models.Token{}, errors.NewAuthorizationError()
	}

	token := new(models.Token)

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(token)

	if err != nil {
		return models.Token{}, err
	}

	return *token, nil
}
