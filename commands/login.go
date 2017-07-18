package commands

import (
	"fmt"
	"net/url"
	"strings"

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
	cfg := config.ReadConfig()
	a := api.NewApi(&cfg)

	if cfg.ApiURL == "" && cmd.ServerUrl == "" {
		return errors.NewNoApiUrlSetError()
	}

	if len(cmd.CaCert) > 0 {
		cfg.CaCert = cmd.CaCert
	}

	serverUrl := cmd.ServerUrl

	if serverUrl != "" {
		if !strings.Contains(serverUrl, "://") {
			serverUrl = "https://" + serverUrl
		}

		credhubInfo, err := a.Target(serverUrl, cfg.CaCert, cmd.SkipTlsValidation)
		if err != nil {
			return err
		}

		parsedUrl, _ := url.Parse(cmd.ServerUrl)
		if parsedUrl.Scheme != "https" {
			warning("Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
				" in transit by third parties. Secure HTTPS API endpoints are recommended.")
		} else {
			if cmd.SkipTlsValidation {
				warning("Warning: The targeted TLS certificate has not been verified for this connection.")
				deprecation("Warning: The --skip-tls-validation flag is deprecated. Please use --ca-cert instead.")
			}
		}
		cfg.ApiURL = parsedUrl.String()
		cfg.InsecureSkipVerify = cmd.SkipTlsValidation
		cfg.AuthURL = credhubInfo.AuthServer.Url
		config.WriteConfig(cfg)
	}

	if cmd.ClientName == "" && cmd.ClientSecret == "" {
		promptForMissingCredentials(&cmd)

	}

	token, err := a.Login(cmd.Username, cmd.Password, cmd.ClientName, cmd.ClientSecret)

	if err != nil {
		a.Logout()
		MarkTokensAsRevokedInConfig(&cfg)
		config.WriteConfig(cfg)
		return err
	} else {
		cfg.AccessToken = token.AccessToken
		cfg.RefreshToken = token.RefreshToken
		config.WriteConfig(cfg)
	}

	if serverUrl != "" {
		fmt.Println("Setting the target url:", cfg.ApiURL)
	}

	fmt.Println("Login Successful")

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
