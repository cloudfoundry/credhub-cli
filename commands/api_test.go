package commands_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
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

	var (
		goodServer *Server
	)

	BeforeEach(func() {
		goodServer = NewServer()

		goodServer.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/info"),
				RespondWith(http.StatusOK, ""),
			),
		)
	})

	AfterEach(func() {
		goodServer.Close()
	})

	It("sets the target URL", func() {
		apiServer := goodServer.URL()
		session := runCommand("api", apiServer)

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))
	})

	It("sets the target URL using a flag", func() {
		apiServer := goodServer.URL()
		session := runCommand("api", "-s", apiServer)

		Eventually(session).Should(Exit(0))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))
	})

	It("will prefer the arguement URL over the flag", func() {
		apiServer := goodServer.URL()
		session := runCommand("api", "-s", "woooo.com", apiServer)

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))

		session = runCommand("api")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(apiServer))
	})

	It("sets the target URL without http", func() {
		session := runCommand("api", goodServer.Addr())

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(goodServer.URL()))

		session = runCommand("api")

		Eventually(session.Out).Should(Say(goodServer.URL()))
	})

	Describe("Validating the target API URL", func() {
		It("fails to set the target when the url is not valid", func() {
			apiServer := goodServer.URL()
			session := runCommand("api", apiServer)

			Eventually(session).Should(Exit(0))

			badServer := NewServer()
			badServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/info"),
					RespondWith(http.StatusNotFound, ""),
				),
			)

			session = runCommand("api", badServer.URL())

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("The targeted API does not appear to be valid."))

			session = runCommand("api")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(goodServer.URL()))

			badServer.Close()
		})
	})
})
