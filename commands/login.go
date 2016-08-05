package commands

import (
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/errors"
)

type LoginCommand struct {
	Username  string `short:"u" long:"username" description:"Sets username"`
	Password  string `short:"p" long:"password" description:"Sets password"`
	ServerUrl string `short:"s" long:"server" description:"API endpoint"`
}

func (cmd LoginCommand) Execute([]string) error {
	cfg, _ := config.ReadConfig()

	if cmd.ServerUrl != "" {
		err := GetApiInfo(&cfg, cmd.ServerUrl)
		if err != nil {
			return err
		}
	}

	err := getUsernameAndPassword(&cmd)
	if err != nil {
		return err
	}

	token, err := actions.NewAuthToken(client.NewHttpClient(cfg.AuthURL), cfg).GetAuthToken(cmd.Username, cmd.Password)
	if err != nil {
		return err
	}

	cfg.AccessToken = token.AccessToken
	cfg.RefreshToken = token.RefreshToken
	config.WriteConfig(cfg)
	fmt.Println("Login Successful")
	return nil
}

func getUsernameAndPassword(cmd *LoginCommand) error {
	if cmd.Username == "" && cmd.Password != "" {
		return errors.NewAuthorizationParametersError()
	}
	if cmd.Username == "" {
		promptForInput("username: ", &cmd.Username)
	}
	if cmd.Password == "" {
		promptForInputWithoutEcho("password: ", &cmd.Password)
	}
	return nil
}

func promptForInput(prompt string, value *string) {
	fmt.Printf(prompt)
	fmt.Scanln(value)
}

func promptForInputWithoutEcho(prompt string, value *string) {
	fmt.Printf(prompt)
	pass, _ := gopass.GetPasswdMasked()
	*value = string(pass)
}
