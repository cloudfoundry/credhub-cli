package api

import (
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

func Api(serverUrl string, caCert []string, skipTlsValidation bool) (models.Info, error) {
	var credhubInfo models.Info
	cfg := config.ReadConfig()

	cfg.CaCert = caCert

	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return credhubInfo, err
	}

	cfg.ApiURL = parsedUrl.String()

	cfg.InsecureSkipVerify = skipTlsValidation

	credhubInfo, err = actions.NewInfo(client.NewHttpClient(cfg), cfg).GetServerInfo()
	if err != nil {
		return credhubInfo, err
	}

	return credhubInfo, nil
}
