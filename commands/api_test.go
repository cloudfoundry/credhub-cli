package commands_test

import (
	"net/http"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/cm-cli/config"
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

	It("revokes existing auth tokens when setting a new api successfully with a different auth server", func() {
		newAuthServer := NewTLSServer()

		apiServer := NewTLSServer()
		apiServer.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/info"),
				RespondWith(http.StatusOK, `{
					"app":{"version":"0.1.0 build DEV","name":"Pivotal Credential Manager"},
					"auth-server":{"url":"`+newAuthServer.URL()+`"}
					}`),
			),
		)

		authServer.AppendHandlers(
			CombineHandlers(
				VerifyRequest("DELETE", "/oauth/token/revoke/5b9c9fd51ba14838ac2e6b222d487106-r"),
				RespondWith(http.StatusOK, ""),
			),
		)

		cfg, _ := config.ReadConfig()
		cfg.AccessToken = "fake_token"
		cfg.RefreshToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImxlZ2FjeS10b2tlbi1rZXkiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiI1YjljOWZkNTFiYTE0ODM4YWMyZTZiMjIyZDQ4NzEwNi1yIiwic3ViIjoiYzE0ZGJjZGQtNzNkOC00ZDdjLWI5NDctYzM4ODVhODAxYzY2Iiwic2NvcGUiOlsiY3JlZGh1Yi53cml0ZSIsImNyZWRodWIucmVhZCJdLCJpYXQiOjE0NzEzMTAwMTIsImV4cCI6MTQ3MTM5NjQxMiwiY2lkIjoiY3JlZGh1YiIsImNsaWVudF9pZCI6ImNyZWRodWIiLCJpc3MiOiJodHRwczovLzUyLjIwNC40OS4xMDc6ODQ0My9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsInJldm9jYWJsZSI6dHJ1ZSwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9uYW1lIjoiY3JlZGh1Yl9jbGkiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX2lkIjoiYzE0ZGJjZGQtNzNkOC00ZDdjLWI5NDctYzM4ODVhODAxYzY2IiwicmV2X3NpZyI6ImQ3MTkyZmUxIiwiYXVkIjpbImNyZWRodWIiXX0.UAp6Ou24f18mdE0XOqG9RLVWZAx3khNHHPeHfuzmcOUYojtILa0_izlGVHhCtNx07f4M9pcRKpo-AijXRw1vSimSTHBeVCDjuuc2nBdznIMhyQSlPpd2stW-WG7Gix82K4gy4oCb1wlTqsK3UKGYoy8JWs6XZqhoZZ6JZM7-Xjj2zag3Q4kgvEBReWC5an_IP6SeCpNt5xWvGdxtTz7ki1WPweUBy0M73ZjRi9_poQT2JmeSIbrePukkfsfCxHG1vM7ApIdzzhdCx6T_KmmMU3xHqhpI_ueLOuvfHjdBinm2atypeTHD83yRRFxhfjRsG1-XguTn-lo_Z2Jis89r5g"
		config.WriteConfig(cfg)

		session := runCommand("api", apiServer.URL())

		Eventually(session).Should(Exit(0))
		newCfg, _ := config.ReadConfig()
		Expect(newCfg.AccessToken).To(Equal("revoked"))
		Expect(newCfg.RefreshToken).To(Equal("revoked"))
		Expect(authServer.ReceivedRequests()).Should(HaveLen(1))
	})

	It("leaves existing auth tokens intact when setting a new api with the same auth server", func() {
		apiServer := NewTLSServer()
		apiServer.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/info"),
				RespondWith(http.StatusOK, `{
					"app":{"version":"my-version","name":"Pivotal Credential Manager"},
					"auth-server":{"url":"`+authServer.URL()+`"}
					}`),
			),
		)

		cfg, _ := config.ReadConfig()
		cfg.AccessToken = "fake_token"
		cfg.RefreshToken = "fake_refresh"
		config.WriteConfig(cfg)

		session := runCommand("api", apiServer.URL())

		Eventually(session).Should(Exit(0))
		newCfg, _ := config.ReadConfig()
		Expect(newCfg.AccessToken).To(Equal("fake_token"))
		Expect(newCfg.RefreshToken).To(Equal("fake_refresh"))
		Expect(authServer.ReceivedRequests()).Should(HaveLen(0))
	})

	It("retains existing tokens when setting the api fails", func() {
		apiServer := NewTLSServer()
		apiServer.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/info"),
				RespondWith(http.StatusNotFound, ""),
			),
		)

		cfg, _ := config.ReadConfig()
		cfg.AuthURL = authServer.URL()
		cfg.AccessToken = "fake_token"
		cfg.RefreshToken = "fake_refresh"
		config.WriteConfig(cfg)

		session := runCommand("api", apiServer.URL())

		Eventually(session).Should(Exit(1))
		newCfg, _ := config.ReadConfig()
		Expect(newCfg.AccessToken).To(Equal("fake_token"))
		Expect(newCfg.RefreshToken).To(Equal("fake_refresh"))
		Expect(authServer.ReceivedRequests()).Should(HaveLen(0))
	})

	Context("when the provided server url's scheme is https", func() {
		var (
			httpsServer       *Server
			apiHttpsServerUrl string
		)

		BeforeEach(func() {
			httpsServer = NewTLSServer()

			apiHttpsServerUrl = httpsServer.URL()

			httpsServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/info"),
					RespondWith(http.StatusOK, `{
					"app":{"version":"0.1.0 build DEV","name":"Pivotal Credential Manager"},
					"auth-server":{"url":"https://example.com"}
					}`),
				),
			)
		})

		AfterEach(func() {
			httpsServer.Close()
		})

		It("sets the target URL", func() {
			session := runCommand("api", apiHttpsServerUrl)

			Eventually(session).Should(Exit(0))

			session = runCommand("api")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(apiHttpsServerUrl))

			config, _ := config.ReadConfig()

			Expect(config.AuthURL).To(Equal("https://example.com"))
		})

		It("sets the target URL using a flag", func() {
			session := runCommand("api", "-s", apiHttpsServerUrl)

			Eventually(session).Should(Exit(0))

			session = runCommand("api")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(apiHttpsServerUrl))
		})

		It("will prefer the command's argument URL over the flag's argument", func() {
			session := runCommand("api", apiHttpsServerUrl, "-s", "woooo.com")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(apiHttpsServerUrl))

			session = runCommand("api")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(apiHttpsServerUrl))
		})

		It("sets the target IP address to an https URL when no URL scheme is provided", func() {
			session := runCommand("api", httpsServer.Addr())

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(httpsServer.URL()))

			session = runCommand("api")

			Eventually(session.Out).Should(Say(httpsServer.URL()))
		})

		Context("when the provided server url is not valid", func() {
			var (
				badServer *Server
			)

			BeforeEach(func() {
				// confirm we have original good server
				session := runCommand("api", apiHttpsServerUrl)

				Eventually(session).Should(Exit(0))

				badServer = NewTLSServer()
				badServer.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/info"),
						RespondWith(http.StatusNotFound, ""),
					),
				)
			})

			AfterEach(func() {
				badServer.Close()
			})

			It("retains previous target when the url is not valid", func() {
				// fail to validate on bad server
				session := runCommand("api", badServer.URL())

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The targeted API does not appear to be valid."))

				// previous value remains
				session = runCommand("api")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say(httpsServer.URL()))
			})
		})

		Context("saving configuration from server", func() {
			It("saves config", func() {
				session := runCommand("api", httpsServer.URL())
				Eventually(session).Should(Exit(0))

				config, error := config.ReadConfig()
				Expect(error).NotTo(HaveOccurred())
				Expect(config.ApiURL).To(Equal(httpsServer.URL()))
				Expect(config.AuthURL).To(Equal("https://example.com"))
			})

			It("sets file permissions so that the configuration is readable and writeable only by the owner", func() {
				configPath := config.ConfigPath()
				os.Remove(configPath)
				session := runCommand("api", httpsServer.URL())
				Eventually(session).Should(Exit(0))

				statResult, _ := os.Stat(configPath)

				Expect(statResult.Mode().String(), "-rw-------")
			})
		})
	})

	Context("when the provided server url's scheme is http", func() {
		var (
			httpServer *Server
		)

		BeforeEach(func() {
			httpServer = NewServer()

			httpServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/info"),
					RespondWith(http.StatusOK, `{
					"app":{"version":"my-version","name":"Pivotal Credential Manager"},
					"auth-server":{"url":"https://example.com"}
					}`),
				),
			)
		})

		AfterEach(func() {
			httpServer.Close()
		})

		It("does not use TLS", func() {
			session := runCommand("api", httpServer.URL())
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(httpServer.URL()))

			session = runCommand("api")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(httpServer.URL()))
		})

		It("prints warning text", func() {
			session := runCommand("api", httpServer.URL())
			Eventually(session).Should(Exit(0))
			Eventually(session).Should(Say("Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
				" in transit by third parties. Secure HTTPS API endpoints are recommended."))
		})
	})
})
