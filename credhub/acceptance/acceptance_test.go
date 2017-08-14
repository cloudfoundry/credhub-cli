package acceptance_test

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

var _ = Describe("CredHub API Acceptance", func() {

	var ch *CredHub

	BeforeEach(func() {
		var err error
		ch, err = New("https://localhost:9000",
			SkipTLSValidation(),
			AuthBuilder(auth.UaaPasswordGrant("credhub_cli", "", "credhub", "password")))

		Expect(err).ToNot(HaveOccurred())
	})

	It("generates a password", func() {
		password, err := ch.GeneratePassword("/example-password", generate.Password{
			Length: 10,
		}, true)

		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(10))
	})

	It("sets a password", func() {
		password, err := ch.SetPassword("/example-password", values.Password("some-password"), true)
		Expect(err).ToNot(HaveOccurred())

		Expect(password.Value).To(BeEquivalentTo("some-password"))
	})

	It("generates a certificate", func() {
		cred, err := ch.GenerateCertificate("/example-certificate", generate.Certificate{
			CommonName: "example.com",
			SelfSign:   true,
		}, true)

		Expect(err).ToNot(HaveOccurred())
		Expect(cred.Name).To(Equal("/example-certificate"))
	})
})
