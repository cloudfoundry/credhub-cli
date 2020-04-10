package commands_test

import (
	"net/http"

	"fmt"

	"code.cloudfoundry.org/credhub-cli/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Regenerate", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("regenerate", "-n", "test-credential")
	ItRequiresAnAPIToBeSet("regenerate", "-n", "test-credential")
	testAutoLogin := []TestAutoLogin{
		{
			method:              "POST",
			responseFixtureFile: "regenerate_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v1/data",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogin, "regenerate", "-n", "test-credential")

	Describe("Regenerating password", func() {
		It("prints the regenerated password secret in yaml format", func() {
			setupRegenerateServer("password", "my-password-stuffs", `"nu-potatoes"`, `{}`)

			session := runCommand("regenerate", "--name", "my-password-stuffs")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-password-stuff"))
			Eventually(session.Out).Should(Say("type: password"))
			Eventually(session.Out).Should(Say("value: <redacted>"))
		})

		It("prints the regenerated password secret in json format", func() {
			setupRegenerateServer("password", "my-password-stuffs", `"nu-potatoes"`, `{}`)

			session := runCommand("regenerate", "--name", "my-password-stuffs", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(fmt.Sprintf(defaultResponseJSON, "password", "my-password-stuffs", `"<redacted>"`, `{}`)))
		})

		It("prints error when server returns an error", func() {
			server.RouteToHandler("POST", "/api/v1/data",
				CombineHandlers(
					VerifyJSON(fmt.Sprintf(regenerateCredentialRequestJson, "my-password-stuffs")),
					RespondWith(http.StatusBadRequest, `{"error":"The password could not be regenerated because the value was statically set. Only generated passwords may be regenerated."}`),
				),
			)

			session := runCommand("regenerate", "--name", "my-password-stuffs")

			Eventually(session).Should(Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("The password could not be regenerated because the value was statically set. Only generated passwords may be regenerated."))
		})

		Describe("with metadata", func() {
			It("regenerates a secret with metadata", func() {
				setupRegenerateServerWithMetadata(
					"password",
					"my-password-stuffs",
					`"nu-potatoes"`,
					`{"some":{"example":"metadata"}, "array":["metadata"]}`,
				)

				session := runCommand("regenerate", "-n", "my-password-stuffs", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

				Eventually(session).Should(Exit(0))
				metadataOutput := `
metadata:
  array:
  - metadata
  some:
    example: metadata`
				Eventually(string(session.Out.Contents())).Should(ContainSubstring(metadataOutput))
			})

			It("errors when metadata is malformed", func() {
				setupRegenerateServerWithMetadata(
					"password",
					"my-password-stuffs",
					`{error}`,
					`"not-valid-json"`)

				session := runCommand("regenerate", "-n", "my-password-stuffs", "--metadata", `"not-valid-json"`)

				Eventually(session).Should(Exit(1))
				Expect(string(session.Err.Contents())).To(ContainSubstring("The argument for --metadata is not a valid json object. Please update and retry your request."))
			})

			It("errors when server does not support metadata", func() {
				setCachedServerVersion("2.5.0")

				session := runCommand("regenerate", "-n", "my-password-stuffs", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

				Eventually(session).Should(Exit(1))
				Expect(string(session.Err.Contents())).To(ContainSubstring("The --metadata flag is not supported for this version of the credhub server (requires >= 2.6.x). Please remove the flag and retry your request."))
			})
		})
	})

	Describe("help", func() {
		ItBehavesLikeHelp("regenerate", "r", func(session *Session) {
			Expect(session.Err).To(Say("regenerate"))
			Expect(session.Err).To(Say("name"))
		})

		It("has short flags", func() {
			Expect(commands.RegenerateCommand{}).To(SatisfyAll(
				commands.HaveFlag("name", "n"),
			))
		})
	})
})

func setupRegenerateServer(keyType, name, value, params string) {
	server.RouteToHandler("POST", "/api/v1/data",
		CombineHandlers(
			VerifyJSON(fmt.Sprintf(regenerateCredentialRequestJson, name)),
			RespondWith(http.StatusOK, fmt.Sprintf(defaultResponseJSON, keyType, name, value, params)),
		),
	)
}
const regenerateCredentialRequestJson = `{"name":"%s", "regenerate":true}`
const regenerateRequestJSONWithMetadata = `{"name":"%s", "regenerate":true, "metadata":%s}`
const regenerateResponseJSONWithMetadata = `{"type":"%s","id":"` + uuid + `","name":"%s","version_created_at":"` + timestamp + `","value":%s,"metadata":%s}`

func setupRegenerateServerWithMetadata(keyType, name, generatedValue, metadata string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(fmt.Sprintf(regenerateRequestJSONWithMetadata, name, metadata)),
			RespondWith(http.StatusOK, fmt.Sprintf(regenerateResponseJSONWithMetadata, keyType, name, generatedValue, metadata)),
		),
	)
}
