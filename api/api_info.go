package api

import (
	"net/url"
	"strings"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

func ApiInfo(serverUrl string, caCert []string, skipTlsValidation bool) error {
	cfg := config.ReadConfig()

	cfg.CaCert = caCert

	if !strings.Contains(serverUrl, "://") {
		serverUrl = "https://" + serverUrl
	}

	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return err
	}

	cfg.ApiURL = parsedUrl.String()

	cfg.InsecureSkipVerify = skipTlsValidation

	credhubInfo, err := actions.NewInfo(client.NewHttpClient(cfg), cfg).GetServerInfo()
	if err != nil {
		return err
	}

	cfg.AuthURL = credhubInfo.AuthServer.Url

	config.WriteConfig(cfg)

	return nil
}
