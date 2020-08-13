package main // import "code.cloudfoundry.org/credhub-cli"

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/cloudfoundry/bosh-utils/errors"

	"code.cloudfoundry.org/credhub-cli/commands"
	"code.cloudfoundry.org/credhub-cli/config"
	"code.cloudfoundry.org/credhub-cli/credhub"
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
			client, err := cfg.Client()
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
		errorType, ok := err.(*flags.Error)
		if ok {
			if errorType.Type.String() == flags.ErrExpectedArgument.String() && runtime.GOOS == "windows" {
				err = errors.WrapError(err, "Flag parsing in windows will interpret any argument with a '/' prefix as an option. Please remove any prepended '/' from flag arguments as it may be causing the following error")
			}
		}

		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
