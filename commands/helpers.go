package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"gopkg.in/yaml.v2"
)

func initializeCredhubClient(cfg config.Config) (*credhub.CredHub, error) {
	var credhubClient *credhub.CredHub

	err = config.ValidateConfig(cfg)
	if err != nil {
		if !clientCredentialsInEnvironment() || config.ValidateConfigApi(cfg) != nil {
			return nil, err
		}
	}

	if clientCredentialsInEnvironment() {
		credhubClient, err = newCredhubClient(&cfg, os.Getenv("CREDHUB_CLIENT"), os.Getenv("CREDHUB_SECRET"), true)
	} else {
		credhubClient, err = newCredhubClient(&cfg, config.AuthClient, config.AuthPassword, false)
	}
	if err != nil {
		return nil, err
	}

	return credhubClient, nil
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
		credhub.AuthURL(cfg.AuthURL),
		credhub.ServerVersion(cfg.ServerVersion))
	return credhubClient, err
}

func clientCredentialsInEnvironment() bool {
	return os.Getenv("CREDHUB_CLIENT") != "" || os.Getenv("CREDHUB_SECRET") != ""
}
