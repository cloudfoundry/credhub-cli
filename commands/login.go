package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/util"
	"github.com/howeyc/gopass"
)

type LoginCommand struct {
	Username          string   `short:"u" long:"username" description:"Authentication username"`
	Password          string   `short:"p" long:"password" description:"Authentication password"`
	ClientName        string   `long:"client-name" description:"Client name for UAA client grant" env:"CREDHUB_CLIENT"`
	ClientSecret      string   `long:"client-secret" description:"Client secret for UAA client grant" env:"CREDHUB_SECRET"`
	ServerUrl         string   `short:"s" long:"server" description:"URI of API server to target" env:"CREDHUB_SERVER"`
	CaCerts           []string `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections" env:"CREDHUB_CA_CERT"`
	SkipTlsValidation bool     `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
	SSO               bool     `long:"sso" description:"Prompt for a one-time passcode to login"`
	SSOPasscode       string   `long:"sso-passcode" description:"One-time passcode"`
}

func (cmd LoginCommand) Execute([]string) error {
	var (
		accessToken  string
		refreshToken string
		err          error
	)
	cfg := config.ReadConfig()

	if cfg.ApiURL == "" && cmd.ServerUrl == "" {
		return errors.NewNoApiUrlSetError()
	}

	if cmd.ServerUrl != "" {
		cfg.InsecureSkipVerify = cmd.SkipTlsValidation

		serverUrl := util.AddDefaultSchemeIfNecessary(cmd.ServerUrl)
		cfg.ApiURL = serverUrl

		err := cfg.UpdateTrustedCAs(cmd.CaCerts)
		if err != nil {
			return err
		}

		credhubInfo, err := GetApiInfo(serverUrl, cfg.CaCerts, cmd.SkipTlsValidation)
		if err != nil {
			return errors.NewNetworkError(err)
		}
		cfg.AuthURL = credhubInfo.AuthServer.URL

		cfg.ServerVersion = credhubInfo.App.Version

		err = verifyAuthServerConnection(cfg, cmd.SkipTlsValidation)
		if err != nil {
			return errors.NewNetworkError(err)
		}
	}

	err = validateParameters(&cmd)

	if err != nil {
		return err
	}
	credhubClient, err := credhub.New(cfg.ApiURL, credhub.CaCerts(cfg.CaCerts...), credhub.SkipTLSValidation(cfg.InsecureSkipVerify))

	uaaClient := uaa.Client{
		AuthURL: cfg.AuthURL,
		Client:  credhubClient.Client(),
	}

	if cmd.ClientName != "" || cmd.ClientSecret != "" {
		accessToken, err = uaaClient.ClientCredentialGrant(cmd.ClientName, cmd.ClientSecret)
	} else {
		err = promptForMissingCredentials(&cmd, &uaaClient)
		if err == nil {
			if cmd.SSOPasscode != "" {
				accessToken, refreshToken, err = uaaClient.PasscodeGrant(config.AuthClient, config.AuthPassword, cmd.SSOPasscode)
			} else {
				accessToken, refreshToken, err = uaaClient.PasswordGrant(config.AuthClient, config.AuthPassword, cmd.Username, cmd.Password)
			}
		}
	}

	if err != nil {
		RevokeTokenIfNecessary(cfg)
		MarkTokensAsRevokedInConfig(&cfg)
		config.WriteConfig(cfg)
		return errors.NewAuthorizationError()
	}

	cfg.AccessToken = accessToken
	cfg.RefreshToken = refreshToken
	if err := config.WriteConfig(cfg); err != nil {
		return err
	}

	if cmd.ServerUrl != "" {
		PrintWarnings(cmd.ServerUrl, cmd.SkipTlsValidation)
		fmt.Println("Setting the target url:", cfg.ApiURL)
	}

	fmt.Println("Login Successful")

	return nil
}

func validateParameters(cmd *LoginCommand) error {
	switch {
	// Intent is client credentials
	case cmd.ClientName != "" || cmd.ClientSecret != "":
		// Make sure nothing else is specified
		if cmd.Username != "" || cmd.Password != "" || cmd.SSO || cmd.SSOPasscode != "" {
			return errors.NewMixedAuthorizationParametersError()
		}

		// Make sure all required fields are specified
		if cmd.ClientName == "" || cmd.ClientSecret == "" {
			return errors.NewClientAuthorizationParametersError()
		}

		return nil

	// Intent is SSO passcode
	case cmd.SSOPasscode != "":
		// Make sure nothing else is specified
		if cmd.ClientName != "" || cmd.ClientSecret != "" || cmd.Username != "" || cmd.Password != "" || cmd.SSO {
			return errors.NewMixedAuthorizationParametersError()
		}

		return nil

	// Intent is to be prompted for token
	case cmd.SSO:
		// Make sure nothing else is specified
		if cmd.ClientName != "" || cmd.ClientSecret != "" || cmd.Username != "" || cmd.Password != "" || cmd.SSOPasscode != "" {
			return errors.NewMixedAuthorizationParametersError()
		}

		return nil

	// Intent is username/password
	default:
		// Make sure nothing else is specified
		if cmd.ClientName != "" || cmd.ClientSecret != "" || cmd.SSO || cmd.SSOPasscode != "" {
			return errors.NewMixedAuthorizationParametersError()
		}

		// Make sure all required fields are specified
		if cmd.Username == "" && cmd.Password != "" {
			return errors.NewPasswordAuthorizationParametersError()
		}

		return nil
	}
}

func promptForMissingCredentials(cmd *LoginCommand, uaa *uaa.Client) error {
	if cmd.SSO || cmd.SSOPasscode != "" {
		if cmd.SSOPasscode == "" {
			md, err := uaa.Metadata()
			if err != nil {
				return err
			}
			fmt.Printf("%s : ", md.PasscodePrompt())
			code, err := gopass.GetPasswdMasked()
			if err != nil {
				return err
			}
			cmd.SSOPasscode = string(code)
		}
		return nil
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
