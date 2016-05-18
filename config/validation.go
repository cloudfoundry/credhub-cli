package config

import "github.com/pivotal-cf/cm-cli/errors"

func ValidateConfig(c Config) error {
	if c.ApiURL == "" {
		return errors.NewNoTargetUrlError()
	}

	return nil
}
