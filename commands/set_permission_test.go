package commands_test

import (
	"code.cloudfoundry.org/credhub-cli/commands"
	"code.cloudfoundry.org/credhub-cli/credhub"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Set Permission", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")
	ItRequiresAnAPIToBeSet("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "get_permission_response.json",
			responseStatus:      http.StatusNotFound,
			endpoint:            "/api/v2/permissions",
		},
		{
			method:              "POST",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")

	Describe("Help", func() {
		ItBehavesLikeHelp("set-permission", "", func(session *Session) {
			Expect(session.Err).To(Say("set-permission"))
			Expect(session.Err).To(Say("actor"))
			Expect(session.Err).To(Say("path"))
		})
	})

	Describe("Set Permission", func() {
		It("parses operation list correctly", func() {
			operationsInput := "read, write, delete"
			parsedInput := commands.ParseOperations(operationsInput)
			expectedOutput := []string{"read", "write", "delete"}
			Expect(parsedInput).To(Equal(expectedOutput))

		})

		It("fails when server version is <2.0", func() {
			ch, _ := credhub.New("https://example.com", credhub.ServerVersion("1.0.0"))
			clientCommand := commands.ClientCommand{}
			clientCommand.SetClient(ch)
			setCommand := commands.SetPermissionCommand{
				Actor:         "some-actor",
				Path:          "'/some-path'",
				Operations:    "read",
				ClientCommand: clientCommand,
			}
			err := setCommand.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Eventually(err).Should(Equal(fmt.Errorf("credhub server version <2.0 not supported")))
		})

		Context("when permission exists", func() {
			Context("when output json flag is used", func() {
				It("updates existing permission", func() {
					responseJsonWithoutDelete := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write"]`)
					responseJsonWithDelete := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write", "delete"]`)
					server.RouteToHandler("GET", "/api/v2/permissions",
						CombineHandlers(
							VerifyRequest("GET", "/api/v2/permissions"),
							RespondWith(http.StatusOK, responseJsonWithoutDelete),
						),
					)

					body := fmt.Sprintf(ADD_PERMISSIONS_REQUEST_JSON, "'/some-path'", "some-actor", `["read", "write", "delete"]`)
					server.RouteToHandler("PUT", "/api/v2/permissions/"+UUID,
						CombineHandlers(
							VerifyRequest("PUT", "/api/v2/permissions/"+UUID),
							VerifyJSON(body),
							RespondWith(http.StatusOK, responseJsonWithDelete),
						),
					)
					session := runCommand("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete", "-j")
					Eventually(session).Should(Exit(0))
					Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "` + UUID + `",
					"actor": "some-actor",
					"path": "'/some-path'",
					"operations": ["read", "write", "delete"]
				}
				`))
				})
			})
			It("updates existing permission", func() {
				responseJsonWithoutDelete := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write"]`)
				responseJsonWithDelete := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write", "delete"]`)
				server.RouteToHandler("GET", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("GET", "/api/v2/permissions"),
						RespondWith(http.StatusOK, responseJsonWithoutDelete),
					),
				)

				body := fmt.Sprintf(ADD_PERMISSIONS_REQUEST_JSON, "'/some-path'", "some-actor", `["read", "write", "delete"]`)
				server.RouteToHandler("PUT", "/api/v2/permissions/"+UUID,
					CombineHandlers(
						VerifyRequest("PUT", "/api/v2/permissions/"+UUID),
						VerifyJSON(body),
						RespondWith(http.StatusOK, responseJsonWithDelete),
					),
				)
				session := runCommand("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("actor: some-actor"))
				Eventually(session.Out).Should(Say(`operations:
- read
- write
- delete`))
				Eventually(session.Out).Should(Say("path: .*'/some-path'.*"))
				Eventually(session.Out).Should(Say("uuid: " + UUID))
			})
		})

		Context("when permission does not exist", func() {
			It("creates a new permission", func() {
				server.RouteToHandler("GET", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("GET", "/api/v2/permissions"),
						RespondWith(http.StatusNotFound, `{"error": "The request could not be completed because the permission does not exist or you do not have sufficient authorization."}`),
					))

				responseJson := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write", "delete"]`)
				body := fmt.Sprintf(ADD_PERMISSIONS_REQUEST_JSON, "'/some-path'", "some-actor", `["read", "write", "delete"]`)

				server.RouteToHandler("POST", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("POST", "/api/v2/permissions"),
						VerifyJSON(body),
						RespondWith(http.StatusOK, responseJson),
					))

				session := runCommand("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("actor: some-actor"))
				Eventually(session.Out).Should(Say(`operations:
- read
- write
- delete`))
				Eventually(session.Out).Should(Say("path: .*'/some-path'.*"))
				Eventually(session.Out).Should(Say("uuid: " + UUID))
			})

			Context("when output json flag is used", func() {
				It("creates a new permission", func() {
					server.RouteToHandler("GET", "/api/v2/permissions",
						CombineHandlers(
							VerifyRequest("GET", "/api/v2/permissions"),
							RespondWith(http.StatusNotFound, `{"error": "The request could not be completed because the permission does not exist or you do not have sufficient authorization."}`),
						))

					responseJson := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write", "delete"]`)
					body := fmt.Sprintf(ADD_PERMISSIONS_REQUEST_JSON, "'/some-path'", "some-actor", `["read", "write", "delete"]`)

					server.RouteToHandler("POST", "/api/v2/permissions",
						CombineHandlers(
							VerifyRequest("POST", "/api/v2/permissions"),
							VerifyJSON(body),
							RespondWith(http.StatusOK, responseJson),
						))

					session := runCommand("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete", "-j")
					Eventually(session).Should(Exit(0))
					Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "` + UUID + `",
					"actor": "some-actor",
					"path": "'/some-path'",
					"operations": ["read", "write", "delete"]
				}
				`))
				})
			})
		})

	})

})
