package acceptance_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-user")
		opts := generate.User{Length: 10}

		By("generate a user with path " + name)
		user, err := credhubClient.GenerateUser(name, opts, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value.Password).To(HaveLen(10))
		generatedUser := user.Value

		By("generate the user again without overwrite returns same user")
		user, err = credhubClient.GenerateUser(name, opts, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value).To(Equal(generatedUser))
	})
})
