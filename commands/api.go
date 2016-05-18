package commands

import (
	"fmt"

	"net/http"
	"net/url"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	. "github.com/pivotal-cf/cm-cli/errors"
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
		parsedUrl, err := url.Parse(serverUrl)
		if err != nil {
			return err
		}
		if parsedUrl.Scheme == "" {
			parsedUrl.Scheme = "http"
		}

		c.ApiURL = parsedUrl.String()

		err = validateTarget(c.ApiURL)
		if err != nil {
			return err
		}
		fmt.Println("Setting the target url:", c.ApiURL)

		config.WriteConfig(c)

	}

	return nil
}

func validateTarget(targetUrl string) error {
	request := client.NewInfoRequest(targetUrl)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return NewNetworkError()
	}

	if response.StatusCode != 200 {
		return NewInvalidTargetError()
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
