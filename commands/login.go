package commands

import (
	"fmt"

	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

type LoginCommand struct {
	Username string `short:"u" long:"username" description:"Sets username"`
	Password string `short:"p" long:"password" description:"Sets password"`
}

func (cmd LoginCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	token, err := actions.NewAuthToken(client.NewHttpClient(cfg.AuthURL), cfg).GetAuthToken(cmd.Username, cmd.Password)

	if err != nil {
		return err
	}

	cfg.AccessToken = token.AccessToken
	config.WriteConfig(cfg)
	fmt.Println("Login Successful")
	return nil
}
