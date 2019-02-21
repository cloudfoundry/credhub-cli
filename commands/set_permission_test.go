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

var _ = Describe("Set Permission", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("set-permission", "-a", "test-actor", "-p", "'/path'", "-o", "read, write, delete")
	ItRequiresAnAPIToBeSet("set-permission", "-a", "test-actor", "-p", "'/path'", "-o", "read, write, delete")

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
	ItAutomaticallyLogsIn(testAutoLogIns, "set-permission", "-a", "test-actor", "-p", "'/path'", "-o", "read, write, delete")

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
				Actor:         "testactor",
				Path:          "'/path'",
				Operations:    "read",
				ClientCommand: clientCommand,
			}
			err := setCommand.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Eventually(err).Should(Equal(fmt.Errorf("credhub server version <2.0 not supported")))
		})

		Context("when permission exists", func() {
			It("updates existing permission", func() {
				server.RouteToHandler("GET", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("GET", "/api/v2/permissions"),
						RespondWith(http.StatusOK, `{"actor": "some-actor",
	"operations": [
   "read",
   "write"
  ],
  "path": "'/some-path'",
  "uuid": "1234"}`),
					),
				)
				server.RouteToHandler("PUT", "/api/v2/permissions/1234",
					CombineHandlers(
						VerifyRequest("PUT", "/api/v2/permissions/1234"),
						VerifyJSON(`{"actor": "some-actor",
					"path": "'/some-path'",
					"operations": ["read", "write", "delete"]}`),
						RespondWith(http.StatusOK, `{"actor": "some-actor",
	"operations": [
   "read",
   "write",
		"delete"
  ],
  "path": "'/some-path'",
  "uuid": "1234"}`),
					),
				)
				session := runCommand("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")
				Eventually(session).Should(Exit(0))
				Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "1234",
					"actor": "some-actor",
					"path": "'/some-path'",
					"operations": ["read", "write", "delete"]
				}
				`))
			})
		})

		Context("when permission does not exist", func() {
			It("creates a new permission", func() {
				server.RouteToHandler("GET", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("GET", "/api/v2/permissions"),
						RespondWith(http.StatusNotFound, `{"error": "The request could not be completed because the permission does not exist or you do not have sufficient authorization."}`),
					))

				responseJson := `{"actor": "some-actor",
	"operations": [
   "read",
   "write",
		"delete"
  ],
  "path": "'/some-path'",
  "uuid": "1234"}`
				body := `{"actor": "some-actor", "path": "'/some-path'", "operations": ["read", "write", "delete"]}`

				server.RouteToHandler("POST", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("POST", "/api/v2/permissions"),
						VerifyJSON(body),
						RespondWith(http.StatusOK, responseJson),
					))

				session := runCommand("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")
				Eventually(session).Should(Exit(0))
				Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "1234",
					"actor": "some-actor",
					"path": "'/some-path'",
					"operations": ["read", "write", "delete"]
				}
				`))
			})
		})

	})

})
