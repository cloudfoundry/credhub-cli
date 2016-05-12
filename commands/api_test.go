package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("API", func() {
	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("api", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("api"))
			Expect(session.Err).To(Say("SERVER_URL"))
		})
	})

	It("sets the target URL", func() {
		apiServer := "http://example.com"
		session := runCommand("api", apiServer)

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))
	})

	It("sets the target URL using a flag", func() {
		apiServer := "http://example.com"
		session := runCommand("api", "-s", apiServer)

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))
	})

	It("will prefer the arguement URL over the flag", func() {
		apiServer := "http://example.com"
		session := runCommand("api", "-s", "woooo.com", apiServer)

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))
	})

	It("sets the target URL without http", func() {
		session := runCommand("api", "example.com")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("http://example.com"))

		session = runCommand("api")

		Eventually(session.Out).Should(Say("http://example.com"))
	})

	It("handles domains that start with http", func() {
		session := runCommand("api", "httpotatoes.com")

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("http://httpotatoes.com"))
	})

	It("handles https URLs", func() {
		session := runCommand("api", "https://example.com")

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out.Contents()).Should(MatchRegexp("^https://example.com"))
	})
})
