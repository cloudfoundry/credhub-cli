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

var _ = Describe("Login", func() {
	Describe("provided a username and password", func() {
		It("authenticates with the UAA server and saves a token", func() {

			uaaServer := NewServer()
			uaaServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token/"),
					VerifyBody([]byte(`grant_type=password&password=pass&response_type=token&username=user`)),
					RespondWith(http.StatusOK, `{
						"access_token":"2YotnFZFEjr1zCsicMWpAA",
						"token_type":"bearer",
						"expires_in":3600}`),
				),
			)

			cfg := config.ReadConfig()
			cfg.AuthURL = uaaServer.URL()
			config.WriteConfig(cfg)

			session := runCommand("login", "-u", "user", "-p", "pass")
			Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("Login Successful"))
			Expect(config.ReadConfig().AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
		})
	})

	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("login", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("login"))
			Expect(session.Err).To(Say("username"))
			Expect(session.Err).To(Say("password"))
		})
	})
})
