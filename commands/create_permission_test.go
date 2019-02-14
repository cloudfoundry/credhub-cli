package commands_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = FDescribe("Create Permission", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("create-permission", "-a", "test-actor", "-p", "/some-path", "-o", "read, write, delete")
	ItRequiresAnAPIToBeSet("create-permission", "-a", "test-actor", "-p", "/some-path", "-o", "read, write, delete")
	ItAutomaticallyLogsIn("POST", "set_response.json", "/api/v2/permissions", "create-permission", "-a", "test-actor", "-p", "/some-path", "-o", "read, write, delete")

	Describe("Help", func() {
		ItBehavesLikeHelp("create-permission", "", func(session *Session) {
			Expect(session.Err).To(Say("create-permission"))
			Expect(session.Err).To(Say("actor"))
			Expect(session.Err).To(Say("path"))
		})
	})

	Describe("Create Permission", func() {
		It("creates a permission", func() {
			responseJson := fmt.Sprintf(ADD_PERMISSIONS_RESPONSE_JSON, "/some-path", "test-actor", "[\"read\", \"write\", \"delete\"]")
			body := `{"actor": "test-actor", "path": "/some-path", "operations": ["read", "write", "delete"]}`

			server.RouteToHandler("POST", "/api/v2/permissions",
				CombineHandlers(
					VerifyRequest("POST", "/api/v2/permissions"),
					VerifyJSON(body),
					RespondWith(http.StatusOK, responseJson),
				),
			)

			session := runCommand("create-permission", "-a", "test-actor", "-p", "/some-path", "-o", "read, write, delete")

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
	})

})
