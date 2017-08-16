package acceptance_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
)

var (
	currentTestNumber = time.Now().UnixNano()
	credhubClient     *credhub.CredHub
)

var _ = BeforeSuite(func() {
	var err error

	credhubClient, err =
		credhub.New("https://localhost:9000",
			credhub.SkipTLSValidation(),
			credhub.AuthBuilder(uaa.PasswordGrantBuilder("credhub_cli", "", "credhub", "password")))

	Expect(err).ToNot(HaveOccurred())
})

func TestCredhub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func testCredentialPath(credentialName string) string {
	return fmt.Sprintf("/acceptance/%v/%v", currentTestNumber, credentialName)
}

func testCredentialPrefix() string {
	return fmt.Sprintf("/acceptance/%v/", currentTestNumber)
}
