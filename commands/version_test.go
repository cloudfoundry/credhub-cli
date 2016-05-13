package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Version", func() {
	It("displays the version", func() {
		session := runCommand("version")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("CLI Version: 0.1.0"))
	})

	It("displays the version with --version", func() {
		session := runCommand("--version")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("CLI Version: 0.1.0"))
	})
})
