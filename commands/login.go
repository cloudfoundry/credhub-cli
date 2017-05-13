package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/howeyc/gopass"
)

type LoginCommand struct {
	Username          string `short:"u" long:"username" description:"Authentication username"`
	Password          string `short:"p" long:"password" description:"Authentication password"`
	ClientName        string `long:"client-name" description:"Client name for UAA client grant [$CREDHUB_CLIENT]"`
	ClientSecret      string `long:"client-secret" description:"Client secret for UAA client grant [$CREDHUB_SECRET]"`
	ServerUrl         string `short:"s" long:"server" description:"URI of API server to target"`
	SkipTlsValidation bool   `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

type loginStrategy interface {
	Login() error
}

type passwordLogin struct {
	cmd    *LoginCommand
	config *config.Config
}

func (cmd LoginCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	strategy := getLoginStrategy(&cmd, &cfg)

	return strategy.Login()
}

func (login passwordLogin) Login() error {
	if login.cmd.ServerUrl != "" {
		err := GetApiInfo(login.config, login.cmd.ServerUrl, login.cmd.SkipTlsValidation)
		if err != nil {
			return err
		}
	}

	err := getUsernameAndPassword(login.cmd)
	if err != nil {
		return err
	}

	token, err := actions.NewAuthToken(client.NewHttpClient(*login.config), *login.config).GetAuthToken(login.cmd.Username, login.cmd.Password)
	if err != nil {
		RevokeTokenIfNecessary(*login.config)
		MarkTokensAsRevokedInConfig(login.config)
		config.WriteConfig(*login.config)
		return err
	}

	login.config.AccessToken = token.AccessToken
	login.config.RefreshToken = token.RefreshToken
	config.WriteConfig(*login.config)

	if login.cmd.ServerUrl != "" {
		fmt.Println("Setting the target url:", login.config.ApiURL)
	}

	fmt.Println("Login Successful")

	return nil
}

func getLoginStrategy(cmd *LoginCommand, config *config.Config) loginStrategy {
	return passwordLogin{
		cmd:    cmd,
		config: config,
	}
}

func getUsernameAndPassword(cmd *LoginCommand) error {
	if cmd.Username == "" && cmd.Password != "" {
		return errors.NewAuthorizationParametersError()
	}
	if cmd.Username == "" {
		fmt.Printf("username: ")
		fmt.Scanln(&cmd.Username)
	}
	if cmd.Password == "" {
		fmt.Printf("password: ")
		pass, _ := gopass.GetPasswdMasked()
		cmd.Password = string(pass)
	}
	return nil
}
