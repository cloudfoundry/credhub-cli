package acceptance_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

var _ = Describe("Password", func() {
	var ch *CredHub

	BeforeEach(func() {
		var err error
		ch, err = New("https://localhost:9000",
			SkipTLSValidation(),
			AuthBuilder(uaa.PasswordGrantBuilder("credhub_cli", "", "credhub", "password")))

		Expect(err).ToNot(HaveOccurred())
	})

	Specify("password lifecycle", func() {

		name := fmt.Sprintf("/acceptance/password/%v", time.Now().UnixNano())
		generatePassword := generate.Password{Length: 10}

		By("generate a password with path " + name)
		password, err := ch.GeneratePassword(name, generatePassword, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(10))
		firstGeneratedPassword := password.Value

		By("generate the password again without overwrite returns same password")
		password, err = ch.GeneratePassword(name, generatePassword, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(firstGeneratedPassword))

		By("setting the password again without overwrite returns same password")
		password, err = ch.SetPassword(name, values.Password("some-password"), false)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(Equal(firstGeneratedPassword))

		By("overwriting the password with generate")
		password, err = ch.GeneratePassword(name, generatePassword, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(HaveLen(10))
		Expect(password.Value).ToNot(Equal(firstGeneratedPassword))

		By("overwriting the password with set")
		password, err = ch.SetPassword(name, values.Password("some-password"), true)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(BeEquivalentTo("some-password"))

		By("getting the password")
		password, err = ch.GetPassword(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(password.Value).To(BeEquivalentTo("some-password"))

		By("deleting the password")
		err = ch.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = ch.GetPassword(name)
		Expect(err).To(HaveOccurred())
	})
})
