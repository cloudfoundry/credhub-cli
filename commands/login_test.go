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

			setConfigAuthUrl(uaaServer.URL())

			session := runCommand("login", "-u", "user", "-p", "pass")

			Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("Login Successful"))
			Expect(config.ReadConfig().AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
		})
	})

	Context("when logging in with server api target", func() {
		var (
			uaaServer *Server
			apiServer *Server
		)

		BeforeEach(func() {
			uaaServer = NewServer()
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

			apiServer = NewServer()
			apiServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/info"),
					RespondWith(http.StatusOK, fmt.Sprintf(`{
					"app":{"version":"0.1.0 build DEV","name":"Pivotal Credential Manager"},
					"auth-server":{"url":"%s","client":"bar"}
					}`, uaaServer.URL())),
				),
			)
		})

		It("sets the target to the server's url and auth server url", func() {
			session := runCommand("login", "-u", "user", "-p", "pass", "-s", apiServer.URL())

			Expect(apiServer.ReceivedRequests()).Should(HaveLen(1))
			Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("Login Successful"))
			Expect(config.ReadConfig().ApiURL).To(Equal(apiServer.URL()))
			Expect(config.ReadConfig().AuthURL).To(Equal(uaaServer.URL()))
			Expect(config.ReadConfig().AuthClient).To(Equal("bar"))
			Expect(config.ReadConfig().AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
		})

		Context("when api server is unavailable", func() {
			var (
				badServer *Server
			)

			BeforeEach(func() {
				badServer = NewServer()
				badServer.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/info"),
						RespondWith(http.StatusBadGateway, nil),
					),
				)
			})

			It("should not login", func() {
				session := runCommand("login", "-u", "user", "-p", "pass", "-s", badServer.URL())

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The targeted API does not appear to be valid. Please validate the API address and retry your request."))
				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
			})

			It("should not override config's existing API URL value", func() {
				cfg := config.ReadConfig()
				cfg.ApiURL = "foo"
				config.WriteConfig(cfg)

				session := runCommand("login", "-u", "user", "-p", "pass", "-s", badServer.URL())

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The targeted API does not appear to be valid. Please validate the API address and retry your request."))
				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Expect(config.ReadConfig().ApiURL).To(Equal("foo"))
			})
		})

		Context("when credentials are invalid", func() {
			var (
				badUaaServer *Server
			)

			BeforeEach(func() {
				badUaaServer = NewServer()
				badUaaServer.AppendHandlers(
					CombineHandlers(
						VerifyRequest("POST", "/oauth/token/"),
						VerifyBody([]byte(`grant_type=password&password=pass&response_type=token&username=user`)),
						RespondWith(http.StatusUnauthorized, `{
						"error":"unauthorized",
						"error_description":"An Authentication object was not found in the SecurityContext"
						}`),
					),
				)
			})

			It("fails to login", func() {
				setConfigAuthUrl(badUaaServer.URL())

				session := runCommand("login", "-u", "user", "-p", "pass")

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The provided username and password combination are incorrect. Please validate your input and retry your request."))
				Expect(badUaaServer.ReceivedRequests()).Should(HaveLen(1))
			})
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

func setConfigAuthUrl(authUrl string) {
	cfg := config.ReadConfig()
	cfg.AuthURL = authUrl
	config.WriteConfig(cfg)
}
