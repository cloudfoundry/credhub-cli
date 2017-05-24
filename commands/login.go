package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/howeyc/gopass"
)

type LoginCommand struct {
	Username          string `short:"u" long:"username" description:"Authentication username"`
	Password          string `short:"p" long:"password" description:"Authentication password"`
	ClientName        string `long:"client-name" description:"Client name for UAA client grant" env:"CREDHUB_CLIENT"`
	ClientSecret      string `long:"client-secret" description:"Client secret for UAA client grant" env:"CREDHUB_SECRET"`
	ServerUrl         string `short:"s" long:"server" description:"URI of API server to target"`
	SkipTlsValidation bool   `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

func (cmd LoginCommand) Execute([]string) error {
	var (
		token models.Token
		err   error
	)
	cfg := config.ReadConfig()

	if cmd.ServerUrl != "" {
		err = GetApiInfo(&cfg, cmd.ServerUrl, cmd.SkipTlsValidation)
		if err != nil {
			return err
		}
	}
	if cmd.ClientName != "" || cmd.ClientSecret != "" {
		if cmd.ClientName == "" || cmd.ClientSecret == "" {
			return errors.NewClientAuthorizationParametersError()
		}

		token, err = actions.NewAuthToken(client.NewHttpClient(cfg), cfg).GetAuthTokenByClientCredential(cmd.ClientName, cmd.ClientSecret)
	} else {
		err = promptForMissingCredentials(&cmd)
		if err != nil {
			return err
		}

		token, err = actions.NewAuthToken(client.NewHttpClient(cfg), cfg).GetAuthTokenByPasswordGrant(cmd.Username, cmd.Password)
	}

	if err != nil {
		RevokeTokenIfNecessary(cfg)
		MarkTokensAsRevokedInConfig(&cfg)
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

func promptForMissingCredentials(cmd *LoginCommand) error {
	if cmd.Username == "" && cmd.Password != "" {
		return errors.NewPasswordAuthorizationParametersError()
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
