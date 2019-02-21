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

	ItRequiresAuthentication("get-permission", "-a", "test-actor", "-p", "'/some-path'")
	ItRequiresAnAPIToBeSet("get-permission", "-a", "test-actor", "-p", "'/some-path'")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "get-permission", "-a", "test-actor", "-p", "'/some-path'")

	Describe("Help", func() {
		ItBehavesLikeHelp("get-permission", "", func(session *Session) {
			Expect(session.Err).To(Say("get-permission"))
			Expect(session.Err).To(Say("actor"))
			Expect(session.Err).To(Say("path"))
		})
	})

	Describe("Get Permission", func() {
		It("fails when server version is <2.0", func() {
			ch, _ := credhub.New("https://example.com", credhub.ServerVersion("1.0.0"))
			clientCommand := commands.ClientCommand{}
			clientCommand.SetClient(ch)
			getCommand := commands.GetPermissionCommand{
				Actor:         "testactor",
				Path:          "'/path'",
				ClientCommand: clientCommand,
			}
			err := getCommand.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Eventually(err).Should(Equal(fmt.Errorf("credhub server version <2.0 not supported")))
		})

		Context("when permission exists", func() {
			It("returns permission", func() {
				server.RouteToHandler("GET", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("GET", "/api/v2/permissions"),
						RespondWith(http.StatusOK, `{"actor": "test-actor",
	"operations": [
   "read",
   "write"
  ],
  "path": "'/some-path'",
  "uuid": "1234"}`),
					))

				session := runCommand("get-permission", "-a", "test-actor", "-p", "'/some-path'")
				Eventually(session).Should(Exit(0))
				Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "1234",
					"actor": "test-actor",
					"path": "'/some-path'",
					"operations": ["read", "write"]
				}
				`))
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
