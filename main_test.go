package main_test

import (
	"fmt"
	"os/exec"

	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("main", func() {
	Context("when no command is provided", func() {
		It("prints help and exits", func() {
			cmd := exec.Command(commandPath)
			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(Exit(1))

			// The command name is printed the same as the credhub cli binary file name
			Expect(session.Err).To(Say(fmt.Sprintf("%s \\[OPTIONS\\] \\[command\\]", commandName)))
		})
	})
	Context("when extra arguments are provided", func() {
		It("prints help and exits", func() {
			cmd := exec.Command(commandPath, "version", "this")
			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("Usage:"))
		})
	})
	Context("when a prepended / is used in the argument", func() {
		BeforeEach(func() {
			if runtime.GOOS != "windows" {
				Skip("only run test on windows")
			}
		})
		It("raises expected error for windows, otherwise does not raise an error", func() {
			cmd := exec.Command(commandPath, "get", "--name", "/foo/bar")
			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("Flag parsing in windows will interpret any argument with a '/' prefix as an option. Please remove any prepended '/' from flag arguments as it may be causing the following error: expected argument for flag `/n, /name', but got option `/foo/bar'"))
		})
	})
})
