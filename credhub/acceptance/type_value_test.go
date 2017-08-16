package acceptance_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Value Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-value")
		cred := values.Value("some string value")
		cred2 := values.Value("another string value")

		By("setting the value for the first time returns same value")
		value, err := credhubClient.SetValue(name, cred, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred))

		By("setting the value again without overwrite returns same value")
		value, err = credhubClient.SetValue(name, cred2, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred))

		By("overwriting the value with set")
		value, err = credhubClient.SetValue(name, cred2, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred2))

		By("getting the value")
		value, err = credhubClient.GetValue(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(value.Value).To(Equal(cred2))

		By("deleting the value")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetValue(name)
		Expect(err).To(HaveOccurred())
	})
})
