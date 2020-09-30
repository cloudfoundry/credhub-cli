package main // import "code.cloudfoundry.org/credhub-cli"

import (
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/cloudfoundry/bosh-utils/errors"

	"code.cloudfoundry.org/credhub-cli/commands"
	"code.cloudfoundry.org/credhub-cli/config"
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/util"
	"github.com/jessevdk/go-flags"
)

type NeedsClient interface {
	SetClient(*credhub.CredHub)
}
type NeedsConfig interface {
	SetConfig(config.Config)
}

func main() {
	debug.SetTraceback("all")
	parser := flags.NewParser(&commands.CredHub, flags.HelpFlag)
	parser.SubcommandsOptional = true
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		if command == nil {
			parser.WriteHelp(os.Stderr)
			os.Exit(1)
		}

		if timeout := parser.FindOptionByLongName("http-timeout").Value().(*time.Duration); timeout != nil {
			_ = os.Setenv("CREDHUB_HTTP_TIMEOUT", timeout.String())
		}

		if cmd, ok := command.(NeedsConfig); ok {
			cmd.SetConfig(config.ReadConfig())
		}

		if cmd, ok := command.(NeedsClient); ok {
			cfg := config.ReadConfig()
			if err := config.ValidateConfig(cfg); err != nil {
				return err
			}
			clientId := cfg.ClientID
			clientSecret := cfg.ClientSecret
			useClientCredentials := true
			if clientId == "" {
				clientId = config.AuthClient
				clientSecret = config.AuthPassword
				useClientCredentials = false
			}
			client, err := credhub.New(cfg.ApiURL,
				credhub.AuthURL(cfg.AuthURL),
				credhub.CaCerts(cfg.CaCerts...),
				credhub.SkipTLSValidation(cfg.InsecureSkipVerify),
				credhub.Auth(auth.Uaa(
					clientId,
					clientSecret,
					"",
					"",
					cfg.AccessToken,
					cfg.RefreshToken,
					useClientCredentials,
				)),
				credhub.ServerVersion(cfg.ServerVersion),
				credhub.SetHttpTimeout(cfg.HttpTimeout),
			)
			if err != nil {
				return err
			}
			cmd.SetClient(client)
		}

		if len(args) != 0 {
			parser.WriteHelp(os.Stderr)
			os.Exit(1)
		}

		return command.Execute(args)
	}

	_, err := parser.Parse()
	if err != nil {
		flagError, ok := err.(*flags.Error)
		if ok {
			errorType := flagError.Type
			if errorType == flags.ErrExpectedArgument && runtime.GOOS == "windows" {
				err = errors.WrapError(err, "Flag parsing in windows will interpret any argument with a '/' prefix as an option. Please remove any prepended '/' from flag arguments as it may be causing the following error")
			} else if errorType == flags.ErrHelp {
				parser.WriteHelp(os.Stderr)
				os.Exit(1)
			}
		}

		util.Error(err.Error())
		os.Exit(1)
	}
}
