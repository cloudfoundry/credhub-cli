package commands_test

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/credhub-cli/commands"
	"code.cloudfoundry.org/credhub-cli/credhub"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Get Permission", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("get-permission", "-a", "some-actor", "-p", "'/some-path'")
	ItRequiresAnAPIToBeSet("get-permission", "-a", "some-actor", "-p", "'/some-path'")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "get-permission", "-a", "some-actor", "-p", "'/some-path'")

	Context("when help flag is used", func() {
		ItBehavesLikeHelp("get-permission", "", func(session *Session) {
			Expect(session.Err).To(Say("get-permission"))
			Expect(session.Err).To(Say("actor"))
			Expect(session.Err).To(Say("path"))
		})
	})

	Context("when server version is < 2.)", func() {
		It("fails", func() {
			ch, _ := credhub.New("https://example.com", credhub.ServerVersion("1.0.0"))
			clientCommand := commands.ClientCommand{}
			clientCommand.SetClient(ch)
			getCommand := commands.GetPermissionCommand{
				Actor:         "some-actor",
				Path:          "'/some-path'",
				ClientCommand: clientCommand,
			}
			err := getCommand.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Eventually(err).Should(Equal(fmt.Errorf("credhub server version <2.0 not supported")))
		})
	})

	Context("when server version is >= 2.0", func() {
		Context("when permission exists", func() {
			Context("when output-json flag is used", func() {
				It("returns permission", func() {
					responseJson := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write"]`)
					server.RouteToHandler("GET", "/api/v2/permissions",
						CombineHandlers(
							VerifyRequest("GET", "/api/v2/permissions"),
							RespondWith(http.StatusOK, responseJson),
						))

					session := runCommand("get-permission", "-a", "some-actor", "-p", "'/some-path'", "-j")
					Eventually(session).Should(Exit(0))
					Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "` + UUID + `",
					"actor": "some-actor",
					"path": "'/some-path'",
					"operations": ["read", "write"]
				}
				`))
				})
			})

			It("returns permission", func() {
				responseJson := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write"]`)
				server.RouteToHandler("GET", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("GET", "/api/v2/permissions"),
						RespondWith(http.StatusOK, responseJson),
					))

				session := runCommand("get-permission", "-a", "some-actor", "-p", "'/some-path'")
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("actor: some-actor"))
				Eventually(session.Out).Should(Say(`operations:
- read
- write`))
				Eventually(session.Out).Should(Say("path: .*'/some-path'.*"))
				Eventually(session.Out).Should(Say("uuid: " + UUID))
			})
		})

		Context("when permission does not exist", func() {
			It("returns error", func() {
				server.RouteToHandler("GET", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("GET", "/api/v2/permissions"),
						RespondWith(http.StatusNotFound, `{"error": "The request could not be completed because the permission does not exist or you do not have sufficient authorization."}`),
					))
				session := runCommand("get-permission", "-a", "test-actor", "-p", "'/some-path'")
				Eventually(session).Should(Exit(1))
			})
		})
	})
})
