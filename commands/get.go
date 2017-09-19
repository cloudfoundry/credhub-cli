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
	Name             string `short:"n" long:"name" description:"Name of the credential to retrieve"`
	Id               string `long:"id" description:"ID of the credential to retrieve"`
	NumberOfVersions int    `long:"versions" description:"Number of versions of the credential to retrieve"`
	OutputJson       bool   `long:"output-json" description:"Return response in JSON format"`
}

func (cmd GetCommand) Execute([]string) error {
	var (
		credential credentials.Credential
		err        error
	)

	cfg := config.ReadConfig()

	var credhubClient *credhub.CredHub

	if clientCredentialsInEnvironment() {
		credhubClient, err = newCredhubClient(&cfg, os.Getenv("CREDHUB_CLIENT"), os.Getenv("CREDHUB_SECRET"), true)
	} else {
		credhubClient, err = newCredhubClient(&cfg, config.AuthClient, config.AuthPassword, false)
	}
	if err != nil {
		return err
	}

	err = config.ValidateConfig(cfg)
	if err != nil {
		if !clientCredentialsInEnvironment() || config.ValidateConfigApi(cfg) != nil{
			return err
		}
	}

	var arrayOfCredentials []credentials.Credential

	if cmd.Name != "" {
		if cmd.NumberOfVersions != 0 {
			arrayOfCredentials, err = credhubClient.GetNVersions(cmd.Name, cmd.NumberOfVersions)
		} else {
			credential, err = credhubClient.GetLatestVersion(cmd.Name)
		}
	} else if cmd.Id != "" {
		credential, err = credhubClient.GetById(cmd.Id)
	} else {
		return errors.NewMissingGetParametersError()
	}

	if err != nil {
		return err
	}

	if arrayOfCredentials != nil {
		output := map[string][]credentials.Credential{
			"versions": arrayOfCredentials,
		}
		printCredential(cmd.OutputJson, output)
	} else {
		printCredential(cmd.OutputJson, credential)
	}

	return nil
}

func printCredential(outputJson bool, v interface{}) {
	if outputJson {
		s, _ := json.MarshalIndent(v, "", "\t")
		fmt.Println(string(s))
	} else {
		s, _ := yaml.Marshal(v)
		fmt.Println(string(s))
	}
}

func newCredhubClient(cfg *config.Config, clientId string, clientSecret string, usingClientCredentials bool) (*credhub.CredHub, error) {
	credhubClient, err := credhub.New(cfg.ApiURL, credhub.CaCerts(cfg.CaCerts...), credhub.SkipTLSValidation(cfg.InsecureSkipVerify), credhub.Auth(auth.Uaa(
		clientId,
		clientSecret,
		"",
		"",
		cfg.AccessToken,
		cfg.RefreshToken,
		usingClientCredentials,
	)),
		credhub.AuthURL(cfg.AuthURL))

	return credhubClient, err
}

func clientCredentialsInEnvironment() bool {
	return os.Getenv("CREDHUB_CLIENT") != "" || os.Getenv("CREDHUB_SECRET") != ""
}
