package commands

import (
	"fmt"

	"net/url"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type ApiCommand struct {
	Server        ApiPositionalArgs `positional-args:"yes"`
	ServerFlagUrl string            `short:"s" long:"server" description:"API endpoint"`
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
		err := GetApiInfo(&cfg, serverUrl)
		if err != nil {
			return err
		}
		config.WriteConfig(cfg)
	}

	return nil
}

func GetApiInfo(cfg *config.Config, serverUrl string) error {
	serverUrl = AddDefaultSchemeIfNecessary(serverUrl)
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return err
	}

	cfg.ApiURL = parsedUrl.String()

	cmInfo, err := actions.NewInfo(client.NewHttpClient(cfg.ApiURL), *cfg).GetServerInfo()
	if err != nil {
		return err
	}
	cfg.AuthURL = cmInfo.AuthServer.Url
	cfg.AuthClient = cmInfo.AuthServer.Client

	if parsedUrl.Scheme != "https" {
		fmt.Println("\033[38;2;255;255;0m" +
			"Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
			" in transit by third parties. Secure HTTPS API endpoints are recommended." +
			"\033[0m")
	}

	fmt.Println("Setting the target url:", cfg.ApiURL)

	return nil
}

func targetUrl(cmd ApiCommand) string {
	if cmd.Server.ServerUrl != "" {
		return cmd.Server.ServerUrl
	} else {
		return cmd.ServerFlagUrl
	}
}
