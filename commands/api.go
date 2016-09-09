package commands

import (
	"fmt"

	"net/url"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
)

type ApiCommand struct {
	Server            ApiPositionalArgs `positional-args:"yes"`
	ServerFlagUrl     string            `short:"s" long:"server" description:"API endpoint"`
	SkipTlsValidation bool              `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

type ApiPositionalArgs struct {
	ServerUrl string `positional-arg-name:"SERVER_URL" description:"The app name"`
}

func (cmd ApiCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	serverUrl := targetUrl(cmd)

	if serverUrl == "" {
		fmt.Println(cfg.ApiURL)
	} else {
		existingCfg := cfg
		err := GetApiInfo(&cfg, serverUrl, cmd.SkipTlsValidation)

		fmt.Println("Setting the target url:", cfg.ApiURL)

		if err != nil {
			return err
		}
		if existingCfg.AuthURL != cfg.AuthURL {
			SendLogoutIfNecessary(existingCfg)
			cfg = RevokedConfig(cfg)
		}
		config.WriteConfig(cfg)
	}

	return nil
}

func GetApiInfo(cfg *config.Config, serverUrl string, skipTlsValidation bool) error {
	serverUrl = AddDefaultSchemeIfNecessary(serverUrl)
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return err
	}

	cfg.ApiURL = parsedUrl.String()

	cfg.InsecureSkipVerify = skipTlsValidation
	cmInfo, err := actions.NewInfo(client.NewHttpClient(*cfg), *cfg).GetServerInfo()
	if err != nil {
		return err
	}
	cfg.AuthURL = cmInfo.AuthServer.Url

	if parsedUrl.Scheme != "https" {
		fmt.Println("\033[38;2;255;255;0m" +
			"Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
			" in transit by third parties. Secure HTTPS API endpoints are recommended." +
			"\033[0m")
	} else {
		if skipTlsValidation {
			fmt.Println("\033[38;2;255;255;0m" +
				"Warning: The targeted TLS certificate has not been verified for this connection." +
				"\033[0m")
		}
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
