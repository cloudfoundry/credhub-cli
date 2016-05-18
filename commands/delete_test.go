package commands_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/cm-cli/config"
)

var _ = Describe("Delete", func() {
	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("delete", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("delete"))
			Expect(session.Err).To(Say("name"))
		})
	})

	It("deletes a secret", func() {
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("DELETE", "/api/v1/data/my-secret"),
				RespondWith(http.StatusOK, ""),
			),
		)

		session := runCommand("delete", "-n", "my-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("Secret successfully deleted"))
	})

	It("deletes a secret and returns 400", func() {
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("DELETE", "/api/v1/data/my-secret"),
				RespondWith(http.StatusBadRequest, ""),
			),
		)

		session := runCommand("delete", "-n", "my-secret")

		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("Unable to perform the request. Please validate your input and retry your request."))
	})

	It("deletes a secret and returns 500", func() {
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("DELETE", "/api/v1/data/my-secret"),
				RespondWith(http.StatusInternalServerError, ""),
			),
		)

		session := runCommand("delete", "-n", "my-secret")

		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("Unable to perform the request. Please validate your input and retry your request."))
	})

	Describe("Errors", func() {
		It("handles no existing secret", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("DELETE", "/api/v1/data/non-existent-secret"),
					RespondWith(http.StatusNotFound, ""),
				),
			)

			session := runCommand("delete", "-n", "non-existent-secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("Secret not found"))
		})

		It("prints an error when API URL is not set", func() {
			cfg := config.ReadConfig()
			cfg.ApiURL = ""
			config.WriteConfig(cfg)

			session := runCommand("delete", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("API location is not set"))
		})

		It("prints an error when the network request fails", func() {
			cfg := config.ReadConfig()
			cfg.ApiURL = "mashed://potatoes"
			config.WriteConfig(cfg)

			session := runCommand("delete", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("No response received for the command"))
		})
	})
})
