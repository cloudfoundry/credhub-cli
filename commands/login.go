package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/howeyc/gopass"
)

type LoginCommand struct {
	Username          string   `short:"u" long:"username" description:"Authentication username"`
	Password          string   `short:"p" long:"password" description:"Authentication password"`
	ClientName        string   `long:"client-name" description:"Client name for UAA client grant" env:"CREDHUB_CLIENT"`
	ClientSecret      string   `long:"client-secret" description:"Client secret for UAA client grant" env:"CREDHUB_SECRET"`
	ServerUrl         string   `short:"s" long:"server" description:"URI of API server to target"`
	CaCert            []string `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections"`
	SkipTlsValidation bool     `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

func (cmd LoginCommand) Execute([]string) error {
	var (
		err error
	)
	cfg := config.ReadConfig()

	a := api.NewApi(&cfg)

	if cfg.ApiURL == "" && cmd.ServerUrl == "" {
		return errors.NewNoApiUrlSetError()
	}

	cfg.ReadTrustedCAs(cmd.CaCert)

	if cmd.ServerUrl != "" {
		err = GetApiInfo(&cfg, cmd.ServerUrl, cmd.SkipTlsValidation)
		if err != nil {
			return err
		}
	}

	err = validateParameters(&cmd)

	if err != nil {
		return err
	}

	if cmd.ClientName != "" || cmd.ClientSecret != "" {
		_, err = a.Login("", "", cmd.ClientName, cmd.ClientSecret)
	} else {
		promptForMissingCredentials(&cmd)
		_, err = a.Login(cmd.Username, cmd.Password, "", "")
	}

	if err != nil {
		a.Logout()
		config.WriteConfig(cfg)
		return err
	}

	config.WriteConfig(cfg)

	if cmd.ServerUrl != "" {
		fmt.Println("Setting the target url:", cfg.ApiURL)
	}

	fmt.Println("Login Successful")

	return nil
}

func validateParameters(cmd *LoginCommand) error {
	if cmd.ClientName != "" || cmd.ClientSecret != "" {
		if cmd.Username != "" || cmd.Password != "" {
			return errors.NewMixedAuthorizationParametersError()
		}

		if cmd.ClientName == "" || cmd.ClientSecret == "" {
			return errors.NewClientAuthorizationParametersError()
		}
	}

	if cmd.Username == "" && cmd.Password != "" {
		return errors.NewPasswordAuthorizationParametersError()
	}

	return nil
}

func promptForMissingCredentials(cmd *LoginCommand) {
	if cmd.Username == "" {
		fmt.Printf("username: ")
		fmt.Scanln(&cmd.Username)
	}
	if cmd.Password == "" {
		fmt.Printf("password: ")
		pass, _ := gopass.GetPasswdMasked()
		cmd.Password = string(pass)
	}
}
