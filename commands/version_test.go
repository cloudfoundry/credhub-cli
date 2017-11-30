package commands_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Version", func() {
	Context("when the user is authenticated", func() {
		BeforeEach(func() {
			login()

			server.RouteToHandler("GET", "/api/v1/data",
				RespondWith(http.StatusOK, `{"data": []}`),
			)
		})

		Context("when the request succeeds", func() {
			BeforeEach(func() {
				responseJson := `{"app":{"name":"CredHub","version":"0.2.0"}}`

				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusOK, responseJson),
				)
			})

			It("displays the version with --version", func() {
				session := runCommand("--version")

				Eventually(session).Should(Exit(0))
				sout := string(session.Out.Contents())
				testVersion(sout)
				Expect(sout).To(ContainSubstring("Server Version: 0.2.0"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusNotFound, ""),
				)
			})

			It("displays the version with --version", func() {
				session := runCommand("--version")

				Eventually(session).Should(Exit(0))
				sout := string(session.Out.Contents())
				testVersion(sout)
				Expect(sout).To(ContainSubstring("Server Version: Not Found"))
			})
		})
	})

	Context("when the user is logged out", func() {
		BeforeEach(func() {
			responseJson := `{"app":{"name":"CredHub","version":"0.2.0"}}`

			server.RouteToHandler("GET", "/info",
				RespondWith(http.StatusOK, responseJson),
			)
		})

		It("returns the CLI version but not the server version", func() {
			session := runCommand("--version")

			Eventually(session).Should(Exit(0))
			sout := string(session.Out.Contents())
			testVersion(sout)
			Expect(sout).To(ContainSubstring("Server Version: Not Found. Have you targeted and authenticated against a CredHub server?"))
		})
	})

	Context("when the config contains invalid tokens", func() {
		BeforeEach(func() {
			responseJson := `{"app":{"name":"CredHub","version":"0.2.0"}}`

			server.RouteToHandler("GET", "/info",
				RespondWith(http.StatusOK, responseJson),
			)

			server.RouteToHandler("GET", "/api/v1/data",
				RespondWith(http.StatusUnauthorized, ""),
			)

			cfg := config.ReadConfig()
			cfg.RefreshToken = "foo"
			cfg.AccessToken = "bar"
			config.WriteConfig(cfg)
		})

		It("returns the CLI version but not the server version", func() {
			session := runCommand("--version")

			Eventually(session).Should(Exit(0))
			sout := string(session.Out.Contents())
			testVersion(sout)
			Expect(sout).To(ContainSubstring("Server Version: Not Found. Have you targeted and authenticated against a CredHub server?"))
		})
	})
})

func testVersion(sout string) {
	Expect(sout).To(ContainSubstring("CLI Version: test-version"))
	Expect(sout).ToNot(ContainSubstring("build DEV"))
}
