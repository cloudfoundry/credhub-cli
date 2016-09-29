package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/credhub-cli/commands"
	"github.com/pivotal-cf/credhub-cli/models"
)

var _ = Describe("Printer Factory", func() {
	It("returns a print string for secret of type value", func() {
		secret := models.NewSecret("my-value", models.SecretBody{
			ContentType: "value",
			Value: "potatoes",
			UpdatedAt: TIMESTAMP,
		})
		Expect(NewPrinterFactory(secret).PrintableSecret()).To(Equal(responseMyValuePotatoes))
	})

	It("returns a print string for secret of type password", func() {
		secret := models.NewSecret("my-password", models.SecretBody{
			ContentType: "password",
			Value: "potatoes",
			UpdatedAt: TIMESTAMP,
		})
		Expect(NewPrinterFactory(secret).PrintableSecret()).To(Equal(responseMyPasswordPotatoes))
	})

	It("returns a print string for secret of type certificate", func() {
		secret := models.NewSecret("my-secret", models.SecretBody{
			ContentType: "certificate",
			Value: models.Certificate{
				Certificate: "my-cert",
				PrivateKey: "my-priv",
				Ca: "my-ca",
			},
			UpdatedAt: TIMESTAMP,
		})
		Expect(NewPrinterFactory(secret).PrintableSecret()).To(Equal(responseMyCertificate))
	})

	It("returns a print string for secret of type ssh", func() {
		secret := models.NewSecret("foo-key", models.SecretBody{
			ContentType: "ssh",
			Value: models.Ssh{
				PublicKey: "some-public-key",
				PrivateKey: "some-private-key",
			},
			UpdatedAt: TIMESTAMP,
		})
		Expect(NewPrinterFactory(secret).PrintableSecret()).To(Equal(responseMySSHFoo))
	})
})
