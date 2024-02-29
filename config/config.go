package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"code.cloudfoundry.org/credhub-cli/util"
)

const AuthClient = "credhub_cli"
const AuthPassword = ""

type Config struct {
	ConfigWithoutSecrets
	ClientID     string
	ClientSecret string
}

func ConfigDir() string {
	return path.Join(userHomeDir(), ".credhub")
}

func ConfigPath() string {
	return path.Join(ConfigDir(), "config.json")
}

func ReadConfig() Config {
	c := Config{}

	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if !os.IsNotExist(err) {
			return c
		}
	}

	json.Unmarshal(data, &c)

	if server, ok := os.LookupEnv("CREDHUB_SERVER"); ok {
		if util.TokenIsPresent(c.AccessToken) {
			util.Warning(
				`WARNING: Two different login methods were detected:
1. A previously run "credhub login" command created a logged-in state
2. CREDHUB_* envrionment variables containing credentials and log-in information are present

This command will now proceed after attempting to log you in using the CREDHUB_* environment variables from method 2, hence ignoring the current logged-in state from method 1.

If you want to get rid of this warning message, you have two options:
a. Run "credhub logout". This will remove the logged-in state created by the "credhub login" command. Subsequent commands will use the environment variables to log you in.
b. Unset the "CREDHUB_SERVER" environment variable. Subsequent commands will use your logged-in state.

`)

		}
		c.ApiURL = util.AddDefaultSchemeIfNecessary(server)
		c.AuthURL = ""
		c.AccessToken = ""
		c.RefreshToken = ""
	}
	if client, ok := os.LookupEnv("CREDHUB_CLIENT"); ok {
		c.ClientID = client
	}
	if clientSecret, ok := os.LookupEnv("CREDHUB_SECRET"); ok {
		c.ClientSecret = clientSecret
	}
	if caCert, ok := os.LookupEnv("CREDHUB_CA_CERT"); ok {
		certs, err := ReadOrGetCaCerts([]string{caCert})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing CA certificates: %+v", err)
			return c
		}
		c.CaCerts = certs
	}
	if timeoutString, ok := os.LookupEnv("CREDHUB_HTTP_TIMEOUT"); ok {
		timeout, err := time.ParseDuration(timeoutString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing HttpTimeout: %+v", err)
			return c
		}
		c.HttpTimeout = &timeout
	}

	return c
}

func WriteConfig(c Config) error {
	err := makeDirectory()
	if err != nil {
		return err
	}

	configWithoutSecrets := ConvertConfigToConfigWithoutSecrets(c)

	data, err := json.Marshal(configWithoutSecrets)
	if err != nil {
		return err
	}

	configPath := ConfigPath()
	return os.WriteFile(configPath, data, 0600)
}

func RemoveConfig() error {
	return os.Remove(ConfigPath())
}

func (cfg *Config) UpdateTrustedCAs(caCerts []string) error {
	var certs []string

	for _, cert := range caCerts {
		certContents, err := util.ReadFileOrStringFromField(cert)
		if err != nil {
			return err
		}
		certs = append(certs, certContents)
	}

	cfg.CaCerts = certs

	return nil
}

func ReadOrGetCaCerts(caCerts []string) ([]string, error) {
	certs := []string{}

	for _, cert := range caCerts {
		certContents, err := util.ReadFileOrStringFromField(cert)
		if err != nil {
			return certs, err
		}
		certs = append(certs, certContents)
	}

	return certs, nil
}
