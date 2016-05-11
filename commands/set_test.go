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

var _ = Describe("Set", func() {
	It("displays help", func() {
		session := runCommand("set", "-h")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("set"))
		Expect(session.Err).To(Say("--name"))
		Expect(session.Err).To(Say("--secret"))
	})

	It("puts a secret", func() {
		requestJson := `{"value":"super-secret-thing"}`
		responseJson := `{"potatoes":"delicious"}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("PUT", "/api/v1/secret/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-s", "super-secret-thing")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(`"potatoes": "delicious"`))
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
		Eventually(session.Err).Should(Say("unsupported protocol scheme"))
	})
})
