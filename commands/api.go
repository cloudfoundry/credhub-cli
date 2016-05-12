package commands

import (
	"fmt"
	"strings"

	"github.com/pivotal-cf/cm-cli/config"
)

type ApiCommand struct {
	Server ApiPositionalArgs `positional-args:"yes"`
	ServerFlagUrl string `short:"s" long:"server" description:"API endpoint"`
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
		fmt.Println("HERERERE", cmd.Server.ServerUrl)
		if strings.HasPrefix(serverUrl, "http://") || strings.HasPrefix(serverUrl, "https://") {
			c.ApiURL = serverUrl
		} else {
			c.ApiURL = "http://" + serverUrl
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