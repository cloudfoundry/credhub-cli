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

	It("generates a password", func() {
		config := Config{
			ApiUrl:             "https://localhost:9000",
			InsecureSkipVerify: true,
		}

		method := auth.UaaPasswordGrant("credhub_cli", "", "credhub", "password")
		ch, err := New(&config, method)

		Expect(err).ToNot(HaveOccurred())

		password, err := ch.GeneratePassword("/example-password", generate.Password{
			Length: 10,
		}, true)

		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(10))
	})

	It("sets a password", func() {
		config := Config{
			ApiUrl:             "https://localhost:9000",
			InsecureSkipVerify: true,
		}
		method := auth.UaaPasswordGrant("credhub_cli", "", "credhub", "password")

		ch, err := New(&config, method)
		Expect(err).ToNot(HaveOccurred())

		password, err := ch.SetPassword("/example-password", values.Password("some-password"), true)
		Expect(err).ToNot(HaveOccurred())

		Expect(password.Value).To(BeEquivalentTo("some-password"))
	})
})
