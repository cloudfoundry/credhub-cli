package commands_test

import (
	"net/http"

	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/credhub-cli/config"
)

var _ = Describe("Delete", func() {
	Describe("Help", func() {
		ItBehavesLikeHelp("delete", "d", func(session *Session) {
			Expect(session.Err).To(Say("delete"))
			Expect(session.Err).To(Say("name"))
		})
	})

	It("deletes a secret", func() {
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("DELETE", "/api/v1/data", "name=my-secret"),
				RespondWith(http.StatusOK, ""),
			),
		)

		session := runCommand("delete", "-n", "my-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("Secret successfully deleted"))
	})

	Describe("Errors", func() {
		It("prints an error when the network request fails", func() {
			cfg := config.ReadConfig()
			cfg.ApiURL = "mashed://potatoes"
			config.WriteConfig(cfg)

			session := runCommand("delete", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("No response received for the command"))
		})

		It("displays missing required parameter", func() {
			session := runCommand("delete")

			Eventually(session).Should(Exit(1))

			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
			}
		})
	})
})
