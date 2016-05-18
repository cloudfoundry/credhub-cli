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

var _ = Describe("Set", func() {
	It("displays help", func() {
		session := runCommand("set", "-h")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("set"))
		Expect(session.Err).To(Say("name"))
		Expect(session.Err).To(Say("secret"))
	})

	It("displays missing -s required parameter", func() {
		session := runCommand("set", "-n", "my-secret")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("the required flag `-s, --secret' was not specified"))
	})

	It("displays missing -n required parameter", func() {
		session := runCommand("set", "-s", "potatoes")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
	})

	It("displays missing -n and -s required parameter", func() {
		session := runCommand("set")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("the required flags `-n, --name' and `-s, --secret' were not specified"))
	})

	It("puts a secret", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)
		requestJson := `{"value":"potatoes"}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("PUT", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-s", "potatoes")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	It("prints an error when the request fails", func() {
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("PUT", "/api/v1/data/my-secret"),
				RespondWith(http.StatusInternalServerError, nil),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-s", "super-secret-thing")

		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("Unable to perform the request"))
	})

	It("prints an error when the request fails 400", func() {
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("PUT", "/api/v1/data/my-secret"),
				RespondWith(http.StatusBadRequest, nil),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-s", "super-secret-thing")

		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("Unable to perform the request"))
	})

	It("prints an error when API URL is not set", func() {
		cfg := config.ReadConfig()
		cfg.ApiURL = ""
		config.WriteConfig(cfg)

		session := runCommand("set", "-n", "me", "-s", "my-secret")

		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("API location is not set"))
	})

	It("prints an error when the network request fails", func() {
		cfg := config.ReadConfig()
		cfg.ApiURL = "mashed://potatoes"
		config.WriteConfig(cfg)

		session := runCommand("set", "-n", "me", "-s", "my-secret")

		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("No response received for the command"))
	})

	It("returns error when the json response is invalid", func() {
		responseJson := `{"name":"my-secret","blah"}`

		requestJson := `{"value":"potatoes"}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("PUT", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-s", "potatoes")

		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("An error occurred when processing the response. Please validate your input and retry your request."))
	})
})
