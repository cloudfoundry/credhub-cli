package commands_test

import (
	"net/http"

	"runtime"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Get", func() {

	ItBehavesLikeHelp("get", "g", func(session *Session) {
		Expect(session.Err).To(Say("Usage"))
		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] get \\[get-OPTIONS\\]"))
		} else {
			Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] get \\[get-OPTIONS\\]"))
		}
	})

	It("displays missing required parameter", func() {
		session := runCommand("get")

		Eventually(session).Should(Exit(1))

		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
		} else {
			Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
		}
	})

	It("gets a string secret", func() {
		responseJson := fmt.Sprintf(STRING_SECRET_ARRAY_RESPONSE_JSON, "value", "my-value", "potatoes")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-value&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-value")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyValuePotatoes))
	})

	It("gets a password secret", func() {
		responseJson := fmt.Sprintf(STRING_SECRET_ARRAY_RESPONSE_JSON, "password", "my-password", "potatoes")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-password&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-password")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPasswordPotatoes))
	})

	It("gets a certificate secret", func() {
		responseJson := fmt.Sprintf(CERTIFICATE_SECRET_ARRAY_RESPONSE_JSON, "my-secret", "my-ca", "my-cert", "my-priv")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-secret&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyCertificate))
	})

	It("gets an rsa secret", func() {
		responseJson := fmt.Sprintf(RSA_SSH_SECRET_ARRAY_RESPONSE_JSON, "rsa", "foo-rsa-key", "some-public-key", "some-private-key")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=foo-rsa-key&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "foo-rsa-key")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyRSAFoo))
	})

	It("can output json", func() {
		responseJson := fmt.Sprintf(STRING_SECRET_ARRAY_RESPONSE_JSON, "password", "my-password", "potatoes")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-password&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-password", "--output-json")

		Eventually(session).Should(Exit(0))
		Eventually(string(session.Out.Contents())).Should(MatchJSON(`{
			"type": "password",
			"name": "my-password",
			"version_created_at": "` + TIMESTAMP + `",
			"value": "potatoes"
		}`))
	})
})
