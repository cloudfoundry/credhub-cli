package commands_test

import (
	"net/http"

	"fmt"

	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Ca-Generate", func() {
	Describe("generating certificate authorities", func() {
		It("posts a valid root CA", func() {
			var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-ca", "my-cert-generated", "my-priv-generated")
			setupPostCaServer("root", "my-ca", "my-cert-generated", "my-priv-generated")

			session := runCommand("ca-generate",
				"-n", "my-ca",
				"-t", "root",
				"--common-name", "my-common-name",
				"--organization", "my-organization",
				"--organization-unit", "my-unit",
				"--locality", "my-locality",
				"--state", "my-state",
				"--country", "my-country",
				"--key-length", "512",
				"--duration", "364")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

	})

	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("ca-generate", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("ca-generate"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("type"))
			Expect(session.Err).To(Say("common-name"))
			Expect(session.Err).To(Say("organization"))
			Expect(session.Err).To(Say("organization-unit"))
			Expect(session.Err).To(Say("locality"))
			Expect(session.Err).To(Say("state"))
			Expect(session.Err).To(Say("country"))
			Expect(session.Err).To(Say("key-length"))
			Expect(session.Err).To(Say("duration"))
		})

		It("displays missing 'n' option as required parameter", func() {
			session := runCommand("ca-generate")

			Eventually(session).Should(Exit(1))
			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
			}
		})

		It("displays the server provided error when an error is received", func() {
			server.AppendHandlers(
				RespondWith(http.StatusBadRequest, `{"error": "you fail."}`),
			)

			session := runCommand("ca-generate", "-n", "my-ca", "-t", "root")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupPostCaServer(caType, name, certificate, priv string) {
	params := `{
	"common_name":"my-common-name",
	"organization":"my-organization",
	"organization_unit":"my-unit",
	"locality":"my-locality",
	"state":"my-state",
	"country":"my-country",
	"key_length": 512,
	"duration": 364
	}`
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", fmt.Sprintf("/api/v1/ca/%s", name)),
			VerifyJSON(fmt.Sprintf(GENERATE_REQUEST_JSON, caType, params)),
			RespondWith(http.StatusOK, fmt.Sprintf(CA_RESPONSE_JSON, caType, certificate, priv)),
		),
	)
}
