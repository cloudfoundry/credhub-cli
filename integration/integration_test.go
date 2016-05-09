package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration", func() {
	var binaryPath string

	BeforeSuite(func() {
		var err error
		binaryPath, err = gexec.Build("github.com/pivotal-cf/cm-cli")

		Expect(err).ToNot(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	Context("when the command is successful", func() {
		It("exits with 0", func() {
			command := exec.Command(binaryPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
		})
	})
})
