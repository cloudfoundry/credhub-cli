// +build !windows

package config_test

import (
	"net/http"

	"code.cloudfoundry.org/credhub-cli/config"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewCredhubClientFromConfig", func() {
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			ConfigWithoutSecrets: config.ConfigWithoutSecrets{
				ApiURL:  "http://api.example.com",
				AuthURL: "http://auth.example.com",
			},
		}
	})

	Context("when client credentials are supplied", func() {
		It("uses the OAuth auth strategy", func() {
			cfg.ClientID = "a-client-id"
			cfg.ClientSecret = "a-client-secret"

			client, err := config.NewCredhubClientFromConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.Auth).To(BeAssignableToTypeOf(&auth.OAuthStrategy{}))
		})
	})

	Context("when non-client credentials are supplied", func() {
		It("uses the OAuth auth strategy", func() {
			client, err := config.NewCredhubClientFromConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.Auth).To(BeAssignableToTypeOf(&auth.OAuthStrategy{}))
		})
	})

	Context("when a client TLS certificate is supplied", func() {
		BeforeEach(func() {
			cfg.ApiURL = "https://api.example.com"
			cfg.ClientCertPath = "../credhub/fixtures/auth-tls-cert.pem"
			cfg.ClientKeyPath = "../credhub/fixtures/auth-tls-key.pem"
		})

		It("does not use the OAuth auth strategy", func() {
			client, err := config.NewCredhubClientFromConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.Auth).NotTo(BeAssignableToTypeOf(&auth.OAuthStrategy{}))
		})

		It("configures the client certificate on HTTPS requests", func() {
			client, err := config.NewCredhubClientFromConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			transport := client.Client().Transport.(*http.Transport)
			tlsConfig := transport.TLSClientConfig
			Expect(len(tlsConfig.Certificates)).To(Equal(1))
		})
	})
})
