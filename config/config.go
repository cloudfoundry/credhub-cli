package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"syscall"
)

const AuthClient = "credhub"
const AuthPassword = ""

type Config struct {
	ApiURL       string
	AuthURL      string
	AccessToken  string
	RefreshToken string
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

func configDir() string {
	return path.Join(userHomeDir(), ".cm")
}

func ConfigPath() string {
	return path.Join(configDir(), "config.json")
}

func ReadConfig() Config {
	c := Config{}

	data, _ := ioutil.ReadFile(ConfigPath())
	json.Unmarshal(data, &c)

	return c
}

func WriteConfig(c Config) {
	os.MkdirAll(configDir(), 0755)

	data, _ := json.Marshal(c)
	configPath := ConfigPath()
	ioutil.WriteFile(configPath, data, 0600)
	syscall.Chmod(configPath, 0600)
}

func RemoveConfig() {
	os.Remove(ConfigPath())
}
