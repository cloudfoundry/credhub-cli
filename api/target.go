package api

import (
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func (a *Api) Target(serverUrl string, caCert []string, skipTlsValidation bool) (models.Info, error) {
	var credhubInfo models.Info

	a.Config.CaCert = caCert

	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return credhubInfo, err
	}

	a.Config.ApiURL = parsedUrl.String()

	a.Config.InsecureSkipVerify = skipTlsValidation

	credhubInfo, err = actions.NewInfo(client.NewHttpClient(*a.Config), *a.Config).GetServerInfo()
	if err != nil {
		return credhubInfo, err
	}

	return credhubInfo, nil
}
