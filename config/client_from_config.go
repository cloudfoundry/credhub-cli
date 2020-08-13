package config

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
)

func NewCredhubClientFromConfig(cfg Config) (*credhub.CredHub, error) {
	var loginOption credhub.Option
	if cfg.ClientCertPath == "" {
		useClientCredentials := true
		clientId := cfg.ClientID
		clientSecret := cfg.ClientSecret
		if clientId == "" {
			useClientCredentials = false
			clientId = AuthClient
			clientSecret = AuthPassword
		}
		loginOption = credhub.Auth(auth.Uaa(
			clientId,
			clientSecret,
			"",
			"",
			cfg.AccessToken,
			cfg.RefreshToken,
			useClientCredentials,
		))
	} else {
		loginOption = credhub.ClientCert(
			cfg.ClientCertPath,
			cfg.ClientKeyPath,
		)
	}

	client, err := credhub.New(cfg.ApiURL,
		credhub.AuthURL(cfg.AuthURL),
		credhub.CaCerts(cfg.CaCerts...),
		credhub.SkipTLSValidation(cfg.InsecureSkipVerify),
		loginOption,
		credhub.ServerVersion(cfg.ServerVersion),
		credhub.SetHttpTimeout(cfg.HttpTimeout),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
