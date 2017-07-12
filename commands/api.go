package commands

import (
	"fmt"
	"strings"

	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/fatih/color"
)

var warning = color.New(color.Bold, color.FgYellow).PrintlnFunc()
var deprecation = color.New(color.Bold, color.FgRed).PrintlnFunc()

type ApiCommand struct {
	Server            ApiPositionalArgs `positional-args:"yes"`
	ServerFlagUrl     string            `short:"s" long:"server" description:"URI of API server to target"`
	CaCert            []string          `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections"`
	SkipTlsValidation bool              `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

type ApiPositionalArgs struct {
	ServerUrl string `positional-arg-name:"SERVER" description:"URI of API server to target"`
}

func (cmd ApiCommand) Execute([]string) error {
	cfg := config.ReadConfig()

	var serverUrl string
	if cmd.Server.ServerUrl != "" {
		serverUrl = cmd.Server.ServerUrl
	} else {
		serverUrl = cmd.ServerFlagUrl
	}

	if serverUrl == "" {
		if cfg.ApiURL != "" {
			fmt.Println(cfg.ApiURL)
			return nil
		} else {
			return errors.NewNoTargetUrlError()
		}
	}

	var err error
	err = api.ApiInfo(serverUrl, cmd.CaCert, cmd.SkipTlsValidation)

	if !strings.Contains(serverUrl, "://") {
		serverUrl = "https://" + serverUrl
	}
	parsedUrl, _ := url.Parse(serverUrl)

	if parsedUrl.Scheme != "https" {
		warning("Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
			" in transit by third parties. Secure HTTPS API endpoints are recommended.")
	} else {
		if cmd.SkipTlsValidation {
			warning("Warning: The targeted TLS certificate has not been verified for this connection.")
			deprecation("Warning: The --skip-tls-validation flag is deprecated. Please use --ca-cert instead.")
		}
	}

	if err != nil {
		return err
	}

	newCfg := config.ReadConfig()
	fmt.Println("Setting the target url:", newCfg.ApiURL)

	if cfg.AuthURL != newCfg.AuthURL {
		api.Logout()
	}

	return nil
}
