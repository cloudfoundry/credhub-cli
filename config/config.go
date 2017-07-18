package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func (cfg *Config) ReadTrustedCAs(caCerts []string) {
	for _, certPath := range caCerts {
		_, err := os.Stat(certPath)
		handleError(err)
		serverCA, err := ioutil.ReadFile(certPath)
		handleError(err)
		cfg.CaCerts = append(cfg.CaCerts, string(serverCA))
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal("Fatal", err)
	}
}
