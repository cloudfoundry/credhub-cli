package commands

import (
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/errors"
)

type LoginCommand struct {
	Username          string `short:"u" long:"username" description:"Authentication username"`
	Password          string `short:"p" long:"password" description:"Authentication password"`
	ServerUrl         string `short:"s" long:"server" description:"URI of API server to target"`
	SkipTlsValidation bool   `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

func (cmd LoginCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	if cmd.ServerUrl != "" {
		err := GetApiInfo(&cfg, cmd.ServerUrl, cmd.SkipTlsValidation)
		if err != nil {
			return err
		}
	}

	err := getUsernameAndPassword(&cmd)
	if err != nil {
		return err
	}

	token, err := actions.NewAuthToken(client.NewHttpClient(cfg), cfg).GetAuthToken(cmd.Username, cmd.Password)
	if err != nil {
		SendLogoutIfNecessary(cfg)
		cfg = RevokedConfig(cfg)
		config.WriteConfig(cfg)
		return err
	}

	cfg.AccessToken = token.AccessToken
	cfg.RefreshToken = token.RefreshToken
	config.WriteConfig(cfg)

	if cmd.ServerUrl != "" {
		fmt.Println("Setting the target url:", cfg.ApiURL)
	}

	fmt.Println("Login Successful")

	return nil
}

func getUsernameAndPassword(cmd *LoginCommand) error {
	if cmd.Username == "" && cmd.Password != "" {
		return errors.NewAuthorizationParametersError()
	}
	if cmd.Username == "" {
		fmt.Printf("username: ")
		username, _ := gopass.GetPasswd()
		cmd.Username = string(username)
	}
	if cmd.Password == "" {
		fmt.Printf("password: ")
		pass, _ := gopass.GetPasswdMasked()
		cmd.Password = string(pass)
	}
	return nil
}
