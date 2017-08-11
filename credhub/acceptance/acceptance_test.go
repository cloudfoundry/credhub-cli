package acceptance_test

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CredHub", func() {

	It("generates a password", func() {
		config := Config{
			ApiUrl:             "https://localhost:9000",
			InsecureSkipVerify: true,
		}
		method := auth.UaaPasswordGrant("credhub_cli", "", "credhub", "password")
		ch := New(&config, method)

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
		ch := New(&config, method)

		password, err := ch.SetPassword("/example-password", values.Password("some-password"), true)
		Expect(err).ToNot(HaveOccurred())

		Expect(password.Value).To(BeEquivalentTo("some-password"))
	})
})
