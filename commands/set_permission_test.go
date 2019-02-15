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

	ItRequiresAuthentication("set-permission", "-a", "test-actor", "-p", "/some-path", "-o", "read, write, delete")
	ItRequiresAnAPIToBeSet("set-permission", "-a", "test-actor", "-p", "/some-path", "-o", "read, write, delete")
	ItAutomaticallyLogsIn("POST", "set_response.json", "/api/v2/permissions", "set-permission", "-a", "test-actor", "-p", "/some-path", "-o", "read, write, delete")

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
		It("creates a new permission", func() {
			responseJson := fmt.Sprintf(ADD_PERMISSIONS_RESPONSE_JSON, "/some-path", "test-actor", "[\"read\", \"write\", \"delete\"]")
			body := `{"actor": "test-actor", "path": "/some-path", "operations": ["read", "write", "delete"]}`

			server.RouteToHandler("POST", "/api/v2/permissions",
				CombineHandlers(
					VerifyRequest("POST", "/api/v2/permissions"),
					VerifyJSON(body),
					RespondWith(http.StatusOK, responseJson),
				),
			)

			session := runCommand("set-permission", "-a", "test-actor", "-p", "/some-path", "-o", "read, write, delete")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"uuid": "5a2edd4f-1686-4c8d-80eb-5daa866f9f86",
					"actor": "test-actor",
					"path": "/some-path",
					"operations": ["read", "write", "delete"]
				}
				`))
		})

		It("fails when server version is <2.0", func() {
			ch, _ := credhub.New("https://example.com", credhub.ServerVersion("1.0.0"))
			clientCommand := commands.ClientCommand{}
			clientCommand.SetClient(ch)
			setCommand := commands.SetPermissionCommand{
				Actor:         "test-actor",
				Path:          "/some-path",
				Operations:    "read",
				ClientCommand: clientCommand,
			}
			err := setCommand.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Eventually(err).Should(Equal(fmt.Errorf("credhub server version <2.0 not supported")))
		})

	})

})
