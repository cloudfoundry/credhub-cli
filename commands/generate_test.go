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

const GENERATE_REQUEST_JSON = `{"type":"value","parameters":%s}`

var _ = Describe("Set", func() {
	It("generates a secret without parameters", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson("{}"))

		session := runCommand("generate", "-n", "my-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret with length", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"length":42}`))

		session := runCommand("generate", "-n", "my-secret", "-l", "42")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret without upper case", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"exclude_upper":true}`))

		session := runCommand("generate", "-n", "my-secret", "--exclude-upper")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret without lower case", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"exclude_lower":true}`))

		session := runCommand("generate", "-n", "my-secret", "--exclude-lower")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret without special characters", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"exclude_special":true}`))

		session := runCommand("generate", "-n", "my-secret", "--exclude-special")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret without numbers", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"exclude_number":true}`))

		session := runCommand("generate", "-n", "my-secret", "--exclude-number")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("generate", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("generate"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("length"))
		})

		It("displays missing 'n' option as required parameters", func() {
			session := runCommand("generate")

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

			session := runCommand("generate", "-n", "my-secret")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupPostServer(name string, value string, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(VALUE_RESPONSE_JSON, value)),
		),
	)
}

func generateRequestJson(params string) string {
	return fmt.Sprintf(GENERATE_REQUEST_JSON, params)
}
