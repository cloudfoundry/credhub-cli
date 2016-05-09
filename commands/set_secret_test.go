package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Set-Secret", func() {
	Describe("Help", func() {
		It("displays help for user", func() {
			session := runCommand("set-secret", "-h")

			Eventually(session).ShouldNot(Exit(0))
			Expect(session.Err).To(Say("set-secret"))
			Expect(session.Err).To(Say("--secret"))
			Expect(session.Err).To(Say("--key"))
		})
	})

	It("puts a secret", func() {
		session := runCommand("set-secret", "-s", "foo-secret", "-k", "key1:value1")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(`"values":`))
		Eventually(session.Out).Should(Say(`"key1": "value1"`))
	})
})