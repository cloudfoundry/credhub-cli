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

const VALUE_REQUEST_JSON = `{"type":"value", "value":"%s"}`

var _ = Describe("Set", func() {
	It("puts a secret using default type", func() {
		setupPutServer("my-secret", "potatoes")

		session := runCommand("set", "-n", "my-secret", "-v", "potatoes")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("puts a secret using explicit value type", func() {
		setupPutServer("my-secret", "potatoes")

		session := runCommand("set", "-n", "my-secret", "-v", "potatoes", "-t", "value")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
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

		It("displays missing 'v' option as required parameters", func() {
			session := runCommand("set", "-n", "my-secret")

			Eventually(session).Should(Exit(1))

			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/v, /value' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-v, --value' was not specified"))
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

func setupPutServer(name string, value string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(fmt.Sprintf(VALUE_REQUEST_JSON, value)),
			RespondWith(http.StatusOK, fmt.Sprintf(RESPONSE_JSON, value)),
		),
	)
}
