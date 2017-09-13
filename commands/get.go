package commands

import (
	"encoding/json"
	"fmt"

	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"gopkg.in/yaml.v2"
)

type GetCommand struct {
	Name       string `short:"n" long:"name" description:"Name of the credential to retrieve"`
	Id         string `long:"id" description:"ID of the credential to retrieve"`
	OutputJson bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd GetCommand) Execute([]string) error {
	var (
		credential credentials.Credential
		err        error
	)

	cfg := config.ReadConfig()

	var credhubClient *credhub.CredHub

	if os.Getenv("CREDHUB_CLIENT") != "" || os.Getenv("CREDHUB_SECRET") != "" {
		credhubClient, err = credhub.New(cfg.ApiURL, credhub.CaCerts(cfg.CaCerts...), credhub.SkipTLSValidation(cfg.InsecureSkipVerify), credhub.Auth(auth.Uaa(
			os.Getenv("CREDHUB_CLIENT"),
			os.Getenv("CREDHUB_SECRET"),
			"",
			"",
			cfg.AccessToken,
			cfg.RefreshToken,
			true,
		)),
			credhub.AuthURL(cfg.AuthURL))
	} else {
		credhubClient, err = credhub.New(cfg.ApiURL, credhub.CaCerts(cfg.CaCerts...), credhub.SkipTLSValidation(cfg.InsecureSkipVerify), credhub.Auth(auth.Uaa(
			config.AuthClient,
			config.AuthPassword,
			"",
			"",
			cfg.AccessToken,
			cfg.RefreshToken,
			false,
		)),
			credhub.AuthURL(cfg.AuthURL))
	}
	if err != nil {
		return err
	}

	err = config.ValidateConfig(cfg)
	if err != nil {
		if os.Getenv("CREDHUB_CLIENT") == "" && os.Getenv("CREDHUB_SECRET") == "" {
			return err
		}
	}

	if cmd.Name != "" {
		credential, err = getLatestVersionWithTokenRefresh(credhubClient, cmd.Name, &cfg)
	} else if cmd.Id != "" {
		credential, err = getByIdWithTokenRefresh(credhubClient, cmd.Id, &cfg)
	} else {
		return errors.NewMissingGetParametersError()
	}

	if err != nil {
		return err
	}

	if cmd.OutputJson {
		s, _ := json.MarshalIndent(credential, "", "\t")
		fmt.Println(string(s))
	} else {
		s, _ := yaml.Marshal(credential)
		fmt.Println(string(s))
	}

	return nil
}

func getLatestVersionWithTokenRefresh(credhubClient *credhub.CredHub, name string, cfg *config.Config) (credential credentials.Credential, err error) {
	credential, err = credhubClient.GetLatestVersion(name)

	return credential, err
}

func getByIdWithTokenRefresh(credhubClient *credhub.CredHub, id string, cfg *config.Config) (credential credentials.Credential, err error) {
	credential, err = credhubClient.GetById(id)
	return credential, err
}
