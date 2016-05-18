package commands_test

import (
	"net/http"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/cm-cli/config"
)

var _ = Describe("Get", func() {
	It("displays help", func() {
		session := runCommand("get", "-h")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("get"))
		Expect(session.Err).To(Say("name"))
	})

	It("gets a secret", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/secret/my-secret"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	Describe("Errors", func() {
		It("handles no existing secret", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/secret/my-secret"),
					RespondWith(http.StatusNotFound, ""),
				),
			)

			session := runCommand("get", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("Secret not found"))
		})

		It("prints an error when API URL is not set", func() {
			cfg := config.ReadConfig()
			cfg.ApiURL = ""
			config.WriteConfig(cfg)

			session := runCommand("get", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("API location is not set"))
		})

		It("prints an error when the network request fails", func() {
			cfg := config.ReadConfig()
			cfg.ApiURL = "mashed://potatoes"
			config.WriteConfig(cfg)

			session := runCommand("get", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("No response received for the command"))
		})

		It("returns error when the response json is invalid", func() {
			responseJson := `{"name":"my-secret","blah"}`

			server.AppendHandlers(
				CombineHandlers(
					//VerifyRequest("GET", "/api/v1/secret/my-secret"),
					RespondWith(http.StatusOK, responseJson),
				),
			)

			session := runCommand("get", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("An error occurred when processing the response. Please validate your input and retry your request."))
		})
	})
})
