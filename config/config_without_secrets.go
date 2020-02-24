package config

import "time"

type ConfigWithoutSecrets struct {
	ApiURL             string
	AuthURL            string
	AccessToken        string
	RefreshToken       string
	InsecureSkipVerify bool
	CaCerts            []string
	ServerVersion      string
	HttpTimeout        *time.Duration
}

func ConvertConfigToConfigWithoutSecrets(config Config) ConfigWithoutSecrets {
	return ConfigWithoutSecrets{
		ApiURL:             config.ApiURL,
		AuthURL:            config.AuthURL,
		AccessToken:        config.AccessToken,
		RefreshToken:       config.RefreshToken,
		InsecureSkipVerify: config.InsecureSkipVerify,
		CaCerts:            config.CaCerts,
		ServerVersion:      config.ServerVersion,
		HttpTimeout:        config.HttpTimeout,
	}
}
