package commands

import (
	"fmt"

	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/util"
	"github.com/fatih/color"
)

var warning = color.New(color.Bold, color.FgYellow).PrintlnFunc()
var deprecation = color.New(color.Bold, color.FgRed).PrintlnFunc()

type ApiCommand struct {
	Server            ApiPositionalArgs `positional-args:"yes" env:"CREDHUB_SERVER"`
	ServerFlagUrl     string            `short:"s" long:"server" description:"URI of API server to target" env:"CREDHUB_SERVER"`
	CaCerts           []string          `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections" env:"CREDHUB_CA_CERT"`
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
			return errors.NewNoTargetUrlError()
		}
	} else {
		existingCfg := cfg

		err := cfg.UpdateTrustedCAs(cmd.CaCerts)
		if err != nil {
			return err
		}

		err = GetApiInfo(&cfg, serverUrl, cmd.SkipTlsValidation)
		if err != nil {
			return err
		}

		fmt.Println("Setting the target url:", cfg.ApiURL)

		if existingCfg.AuthURL != cfg.AuthURL {
			RevokeTokenIfNecessary(existingCfg)
			MarkTokensAsRevokedInConfig(&cfg)
		}
		err = config.WriteConfig(cfg)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetApiInfo(cfg *config.Config, serverUrl string, skipTlsValidation bool) error {
	serverUrl = util.AddDefaultSchemeIfNecessary(serverUrl)
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return err
	}

	cfg.ApiURL = parsedUrl.String()

	cfg.InsecureSkipVerify = skipTlsValidation
	credhubInfo, err := actions.NewInfo(client.NewHttpClient(*cfg), *cfg).GetServerInfo()
	if err != nil {
		return err
	}
	cfg.AuthURL = credhubInfo.AuthServer.Url

	if parsedUrl.Scheme != "https" {
		warning("Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
			" in transit by third parties. Secure HTTPS API endpoints are recommended.")
	} else {
		if skipTlsValidation {
			warning("Warning: The targeted TLS certificate has not been verified for this connection.")
			deprecation("Warning: The --skip-tls-validation flag is deprecated. Please use --ca-cert instead.")
		}
	}

	err = verifyAuthServerConnection(*cfg, skipTlsValidation)
	if err != nil {
		return err
	}

	return nil
}

func targetUrl(cmd ApiCommand) string {
	if cmd.Server.ServerUrl != "" {
		return cmd.Server.ServerUrl
	} else {
		return cmd.ServerFlagUrl
	}
}

func verifyAuthServerConnection(cfg config.Config, skipTlsValidation bool) error {
	var err error

	if !skipTlsValidation {
		err = actions.VerifyAuthServerConnection(client.NewHttpClient(cfg), cfg)
	}

	return err
}
