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

var _ = Describe("Delete Permission", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("delete-permission", "-a", "some-actor", "-p", "'/some-path'")
	ItRequiresAnAPIToBeSet("delete-permission", "-a", "some-actor", "-p", "'/some-path'")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions",
		},
		{
			method:              "DELETE",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions/" + UUID,
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "delete-permission", "-a", "some-actor", "-p", "'/some-path'")

	Describe("Help", func() {
		ItBehavesLikeHelp("delete-permission", "", func(session *Session) {
			Expect(session.Err).To(Say("delete-permission"))
			Expect(session.Err).To(Say("actor"))
			Expect(session.Err).To(Say("path"))
		})
	})

	Describe("Delete Permission", func() {

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
			It("deletes existing permission", func() {
				responseJson := fmt.Sprintf(PERMISSIONS_RESPONSE_JSON, "'/some-path'", "some-actor", `["read", "write"]`)
				server.RouteToHandler("GET", "/api/v2/permissions",
					CombineHandlers(
						VerifyRequest("GET", "/api/v2/permissions"),
						RespondWith(http.StatusOK, responseJson),
					),
				)
				server.RouteToHandler("DELETE", "/api/v2/permissions/"+UUID,
					CombineHandlers(
						VerifyRequest("DELETE", "/api/v2/permissions/"+UUID),
						RespondWith(http.StatusOK, responseJson),
					),
				)
				session := runCommand("delete-permission", "-a", "some-actor", "-p", "'/some-path'")
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("actor: some-actor"))
				Eventually(session.Out).Should(Say(`operations:
- read
- write`))
				Eventually(session.Out).Should(Say("path: .*'/some-path'.*"))
				Eventually(session.Out).Should(Say("uuid: " + UUID))
			})
		})
	})

	Context("when permission does not exist", func() {
		It("throws an error", func() {
			server.RouteToHandler("GET", "/api/v2/permissions",
				CombineHandlers(
					VerifyRequest("GET", "/api/v2/permissions"),
					RespondWith(http.StatusNotFound, `{"error": "The request could not be completed because the permission does not exist or you do not have sufficient authorization."}`),
				))

			errorJson := `{"error": "The request includes a permission that does not exist."}`
			server.RouteToHandler("DELETE", "/api/v2/permissions",
				CombineHandlers(
					VerifyRequest("DELETE", "/api/v2/permissions"),
					RespondWith(http.StatusOK, errorJson),
				))

			session := runCommand("delete-permission", "-a", "some-actor", "-p", "'/some-path'")
			Eventually(session).Should(Exit(1))

		})
	})

})
