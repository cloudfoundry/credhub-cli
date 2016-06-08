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

var _ = Describe("Set", func() {
	It("puts a secret using default type", func() {
		setupPutValueServer("my-secret", "potatoes")

		session := runCommand("set", "-n", "my-secret", "-v", "potatoes")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("puts a secret using explicit value type", func() {
		setupPutValueServer("my-secret", "potatoes")

		session := runCommand("set", "-n", "my-secret", "-v", "potatoes", "-t", "value")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("puts a secret using explicit certificate type", func() {
		var responseMyCertificate = fmt.Sprintf(CERTIFICATE_RESPONSE_TABLE, "my-secret", "my-ca", "my-pub", "my-priv")
		setupPutCertificateServer("my-secret", "my-ca", "my-pub", "my-priv")

		session := runCommand("set", "-n", "my-secret",
			"-t", "certificate", "--ca-string", "my-ca",
			"--public-string", "my-pub", "--private-string", "my-priv")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyCertificate))
	})

	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("set", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("set"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("secret"))
		})

		It("displays missing 'n' option as required parameter", func() {
			session := runCommand("set", "-v", "potatoes")

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

			session := runCommand("set", "-n", "my-secret", "-v", "tomatoes")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupPutValueServer(name string, value string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(fmt.Sprintf(VALUE_REQUEST_JSON, value)),
			RespondWith(http.StatusOK, fmt.Sprintf(VALUE_RESPONSE_JSON, value)),
		),
	)
}

func setupPutCertificateServer(name string, ca string, pub string, priv string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(fmt.Sprintf(CERTIFICATE_REQUEST_JSON, ca, pub, priv)),
			RespondWith(http.StatusOK, fmt.Sprintf(CERTIFICATE_RESPONSE_JSON, ca, pub, priv)),
		),
	)
}
