package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

var _ = Describe("Find", func() {

	passwordName1 := testCredentialPath("first-password")
	passwordName2 := testCredentialPath("second-password")

	var expectedPassword1 credentials.Password
	var expectedPassword2 credentials.Password

	BeforeEach(func() {
		var err error

		generatePassword := generate.Password{Length: 10}

		expectedPassword1, err = credhubClient.GeneratePassword(passwordName1, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())

		expectedPassword2, err = credhubClient.GeneratePassword(passwordName2, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := credhubClient.Delete(passwordName1)
		Expect(err).ToNot(HaveOccurred())
		err = credhubClient.Delete(passwordName2)
		Expect(err).ToNot(HaveOccurred())
	})

	Specify("finding the credentials by path", func() {
		results, err := credhubClient.FindByPath(testCredentialPrefix())

		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(ConsistOf(expectedPassword1.Base, expectedPassword2.Base))
	})
})
