package commands_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Version", func() {
	Context("when the request succeeds", func() {
		BeforeEach(func() {
			responseJson := `{"app":{"name":"Pivotal Credential Manager","version":"0.2.0"}}`

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/info"),
					RespondWith(http.StatusOK, responseJson),
				),
			)
		})

		It("displays the version with --version", func() {
			session := runCommand("--version")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("CLI Version: 0.1.0"))
			Eventually(session.Out).Should(Say("CM Version: 0.2.0"))
		})

		It("displays the version with -v", func() {
			session := runCommand("-v")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("CLI Version: 0.1.0"))
			Eventually(session.Out).Should(Say("CM Version: 0.2.0"))
		})
	})
})
