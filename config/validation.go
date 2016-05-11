package config

import "errors"

func ValidateConfig(c Config) error {
	if c.ApiURL == "" {
		return errors.New("Credential Manager API location is not set. Please target the location of your credential manager with `cm api` to continue.")
	}

	return nil
}
