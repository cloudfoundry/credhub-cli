package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

var _ = Describe("Password Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-password")
		generatePassword := generate.Password{Length: 10}

		By("generate a password with path " + name)
		password, err := credhubClient.GeneratePassword(name, generatePassword, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(10))
		firstGeneratedPassword := password.Value

		By("generate the password again without overwrite returns same password")
		password, err = credhubClient.GeneratePassword(name, generatePassword, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(firstGeneratedPassword))

		By("setting the password again without overwrite returns same password")
		password, err = credhubClient.SetPassword(name, values.Password("some-password"), false)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(firstGeneratedPassword))

		By("overwriting the password with generate")
		password, err = credhubClient.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(10))
		Expect(password.Value).ToNot(Equal(firstGeneratedPassword))

		By("overwriting the password with set")
		password, err = credhubClient.SetPassword(name, values.Password("some-password"), true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(values.Password("some-password")))

		By("getting the password")
		password, err = credhubClient.GetPassword(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(values.Password("some-password")))

		By("deleting the password")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetPassword(name)
		Expect(err).To(HaveOccurred())
	})
})
