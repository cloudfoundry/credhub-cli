package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Set", func() {
	Describe("Help", func() {
		It("displays help for user", func() {
			session := runCommand("set", "-h")

			Eventually(session).ShouldNot(Exit(0))
			Expect(session.Err).To(Say("set"))
			Expect(session.Err).To(Say("--identifier"))
			Expect(session.Err).To(Say("--secret"))
		})
	})

	It("puts a secret", func() {
		session := runCommand("set", "-i", "my-secret", "-s", "super secret thing")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(`"value": "super secret thing"`))
	})
})