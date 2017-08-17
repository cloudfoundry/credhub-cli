package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
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

		newUser := values.User{Username: "name", Password: "password"}

		By("setting the user again without overwrite returns same user")
		user, err = credhubClient.SetUser(name, newUser, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value).To(Equal(generatedUser))

		By("overwriting the user with generate")
		user, err = credhubClient.GenerateUser(name, opts, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value.Password).To(HaveLen(10))
		Expect(user.Value).ToNot(Equal(generatedUser))

		By("overwriting the user with set")
		user, err = credhubClient.SetUser(name, newUser, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value.User).To(Equal(newUser))

		By("getting the user")
		user, err = credhubClient.GetUser(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Value.User).To(Equal(newUser))

		By("deleting the user")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetUser(name)
		Expect(err).To(HaveOccurred())
	})
})
