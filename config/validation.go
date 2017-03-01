package config

import "github.com/cloudfoundry-incubator/credhub-cli/errors"

func ValidateConfig(c Config) error {
	if c.ApiURL == "" {
		return errors.NewNoTargetUrlError()
	}

	return nil
}
