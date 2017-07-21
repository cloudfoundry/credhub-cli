package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const AuthClient = "credhub_cli"
const AuthPassword = ""

type Config struct {
	ApiURL             string
	AuthURL            string
	AccessToken        string
	RefreshToken       string
	InsecureSkipVerify bool
	CaCerts            []string
}

func ConfigDir() string {
	return path.Join(userHomeDir(), ".credhub")
}

func ConfigPath() string {
	return path.Join(ConfigDir(), "config.json")
}

func ReadConfig() Config {
	c := Config{}

	data, err := ioutil.ReadFile(ConfigPath())
	if err != nil {
		return c
	}

	json.Unmarshal(data, &c)

	return c
}

func WriteConfig(c Config) error {
	err := makeDirectory()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	configPath := ConfigPath()
	return ioutil.WriteFile(configPath, data, 0600)
}

func RemoveConfig() error {
	return os.Remove(ConfigPath())
}

func (cfg *Config) UpdateTrustedCAs(caCerts []string) error {
	certs := []string{}

	for _, cert := range caCerts {
		_, err := os.Stat(cert)

		if err != nil {
			certs = append(certs, string(cert))
		} else {
			certContents, err := ioutil.ReadFile(cert)

			if err != nil {
				return err
			}

			certs = append(certs, string(certContents))
		}
	}

	cfg.CaCerts = certs

	return nil
}
