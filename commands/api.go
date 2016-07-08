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
	c := config.ReadConfig()
	serverUrl := targetUrl(cmd)

	if serverUrl == "" {
		fmt.Println(c.ApiURL)
	} else {
		serverUrl = AddDefaultSchemeIfNecessary(serverUrl)
		parsedUrl, err := url.Parse(serverUrl)
		if err != nil {
			return err
		}

		c.ApiURL = parsedUrl.String()

		action := actions.NewApi(client.NewHttpClient(c))

		err = action.ValidateTarget(c.ApiURL)
		if err != nil {
			return err
		}
		if parsedUrl.Scheme != "https" {
			fmt.Println("Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
				" in transit by third parties. Secure HTTPS API endpoints are recommended.")
		}

		fmt.Println("Setting the target url:", c.ApiURL)

		config.WriteConfig(c)
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
