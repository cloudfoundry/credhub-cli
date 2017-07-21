// +build !windows

package config_test

import (
	"io/ioutil"

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
			ca1, err := ioutil.ReadFile("../test/server-tls-ca.pem")
			Expect(err).To(BeNil())
			ca2, err := ioutil.ReadFile("../test/auth-tls-ca.pem")
			Expect(err).To(BeNil())

			cfg.UpdateTrustedCAs([]string{"../test/server-tls-ca.pem", "../test/auth-tls-ca.pem"})

			Expect(cfg.CaCerts).To(ConsistOf([]string{string(ca1), string(ca2)}))
		})

		It("overrides previous CAs", func() {
			testCa, err := ioutil.ReadFile("../test/server-tls-ca.pem")
			Expect(err).To(BeNil())

			cfg.CaCerts = []string{"cert1", "cert2"}
			cfg.UpdateTrustedCAs([]string{"../test/server-tls-ca.pem"})

			Expect(cfg.CaCerts).To(ConsistOf([]string{string(testCa)}))
		})
	})
})
