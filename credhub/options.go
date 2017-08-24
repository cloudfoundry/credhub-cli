package credhub

import (
	"crypto/x509"
	"errors"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

type Option func(*CredHub) error

func Auth(method auth.Builder) Option {
	return func(c *CredHub) error {
		c.authBuilder = method
		return nil
	}
}

func AuthStrategy(strategy auth.Strategy) Option {
	return func(c *CredHub) error {
		c.Auth = strategy
		return nil
	}
}

func AuthURL(authURL string) Option {
	return func(c *CredHub) error {
		var err error
		c.authURL, err = url.Parse(authURL)
		return err
	}
}

func CaCerts(certs ...string) Option {
	return func(c *CredHub) error {
		c.caCerts = x509.NewCertPool()

		for _, cert := range certs {
			ok := c.caCerts.AppendCertsFromPEM([]byte(cert))
			if !ok {
				return errors.New("provided ca certs are invalid")
			}
		}

		return nil
	}
}

func SkipTLSValidation() Option {
	return func(c *CredHub) error {
		c.insecureSkipVerify = true
		return nil
	}
}
