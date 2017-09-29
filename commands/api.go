package commands

import (
	"fmt"

	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/util"
	"github.com/fatih/color"
)

var warning = color.New(color.Bold, color.FgYellow).PrintlnFunc()
var deprecation = color.New(color.Bold, color.FgRed).PrintlnFunc()

type ApiCommand struct {
	Server            ApiPositionalArgs `positional-args:"yes" env:"CREDHUB_SERVER"`
	ServerFlagUrl     string            `short:"s" long:"server" description:"URI of API server to target" env:"CREDHUB_SERVER"`
	CaCerts           []string          `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections. Multiple flags may be provided." env:"CREDHUB_CA_CERT"`
	SkipTlsValidation bool              `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

type ApiPositionalArgs struct {
	ServerUrl string `positional-arg-name:"SERVER" description:"URI of API server to target"`
}

func (cmd ApiCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	serverUrl := targetUrl(cmd)

	if serverUrl == "" {
		if cfg.ApiURL != "" {
			fmt.Println(cfg.ApiURL)
		} else {
			return errors.NewNoApiUrlSetError()
		}
	} else {
		caCerts, err := ReadOrGetCaCerts(cmd.CaCerts)
		if err != nil {
			return err
		}

		serverUrl = util.AddDefaultSchemeIfNecessary(serverUrl)

		credhubInfo, err := GetApiInfo(serverUrl, caCerts, cmd.SkipTlsValidation)
		if err != nil {
			return errors.NewNetworkError(err)
		}

		if credhubInfo.AuthServer.URL != cfg.AuthURL {
			RevokeTokenIfNecessary(cfg)
			MarkTokensAsRevokedInConfig(&cfg)
		}

		cfg.ApiURL = serverUrl
		cfg.AuthURL = credhubInfo.AuthServer.URL
		cfg.ServerVersion = credhubInfo.App.Version
		cfg.InsecureSkipVerify = cmd.SkipTlsValidation
		cfg.CaCerts = caCerts

		err = VerifyAuthServerConnection(cfg, cmd.SkipTlsValidation)
		if err != nil {
			return errors.NewNetworkError(err)
		}

		err = PrintWarnings(serverUrl, cmd.SkipTlsValidation)
		if err != nil {
			return err
		}
		fmt.Println("Setting the target url:", cfg.ApiURL)

		err = config.WriteConfig(cfg)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetApiInfo(serverUrl string, caCerts []string, skipTlsValidation bool) (*server.Info, error) {
	credhubClient, err := credhub.New(serverUrl, credhub.CaCerts(caCerts...), credhub.SkipTLSValidation(skipTlsValidation))
	if err != nil {
		return nil, err
	}

	credhubInfo, err := credhubClient.Info()
	return credhubInfo, err
}

func PrintWarnings(serverUrl string, skipTlsValidation bool) error {
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return err
	}

	if parsedUrl.Scheme != "https" {
		warning("Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
			" in transit by third parties. Secure HTTPS API endpoints are recommended.")
	} else {
		if skipTlsValidation {
			warning("Warning: The targeted TLS certificate has not been verified for this connection.")
			deprecation("Warning: The --skip-tls-validation flag is deprecated. Please use --ca-cert instead.")
		}
	}
	return nil
}

func ReadOrGetCaCerts(caCerts []string) ([]string, error) {
	certs := []string{}

	for _, cert := range caCerts {
		certContents, err := util.ReadFileOrStringFromField(cert)
		if err != nil {
			return certs, err
		}
		certs = append(certs, certContents)
	}

	return certs, nil
}

func targetUrl(cmd ApiCommand) string {
	if cmd.Server.ServerUrl != "" {
		return cmd.Server.ServerUrl
	} else {
		return cmd.ServerFlagUrl
	}
}

func VerifyAuthServerConnection(cfg config.Config, skipTlsValidation bool) error {
	var err error

	if !skipTlsValidation {
		err = actions.VerifyAuthServerConnection(client.NewHttpClient(cfg), cfg)
	}

	return err
}
