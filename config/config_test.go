// +build !windows

package config_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "http://api.example.com",
			AuthURL: "http://auth.example.com",
		}
	})

	It("places the config file in .cm in the home directory", func() {
		Expect(config.ConfigPath()).To(HaveSuffix(`/.credhub/config.json`))
	})

	Describe("#UpdateTrustedCAs", func() {
		It("reads multiple certs", func() {
			cfg.UpdateTrustedCAs([]string{"../test/server-tls-ca.pem", "../test/auth-tls-ca.pem"})
			Expect(cfg.CaCerts).To(HaveLen(2))
		})

		It("overrides previous CAs", func() {
			cfg.CaCerts = []string{"cert1", "cert2"}

			cfg.UpdateTrustedCAs([]string{"../test/server-tls-ca.pem", "../test/auth-tls-ca.pem"})
			Expect(cfg.CaCerts).To(HaveLen(2))
		})
	})
})
