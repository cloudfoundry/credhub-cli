package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("API", func() {
	It("displays help", func() {
		session := runCommand("api", "-h")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("api"))
		Expect(session.Err).To(Say("--server"))
	})

	It("sets the target URL", func() {
		apiServer := "http://example.com"
		session := runCommand("api", "-s", apiServer)

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))
	})

	It("sets the target URL without http", func() {
		session := runCommand("api", "-s", "example.com")

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("http://example.com"))
	})

	It("handles domains that start with http", func() {
		session := runCommand("api", "-s", "httpotatoes.com")

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("http://httpotatoes.com"))
	})

	It("handles https URLs", func() {
		session := runCommand("api", "-s", "https://example.com")

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out.Contents()).Should(MatchRegexp("^https://example.com"))
	})
})
