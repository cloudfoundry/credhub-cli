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
const GENERATE_REQUEST_JSON = `{"type":"value","parameters":%s}`
const RESPONSE_JSON = `{"value":"%s","type":"value"}`
const RESPONSE_TABLE = `Type:	value\nName:	%s\nValue:	%s`

var responseMyPotatoes = fmt.Sprintf(RESPONSE_TABLE, "my-secret", "potatoes")

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

	It("generates a secret without parameters", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson("{}"))

		session := runCommand("set", "-n", "my-secret", "-g")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret with length", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"length":42}`))

		session := runCommand("set", "-n", "my-secret", "-g", "-l", "42")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret without upper case", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"exclude_upper":true}`))

		session := runCommand("set", "-n", "my-secret", "-g", "--exclude-upper")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret without lower case", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"exclude_lower":true}`))

		session := runCommand("set", "-n", "my-secret", "-g", "--exclude-lower")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret without special characters", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"exclude_special":true}`))

		session := runCommand("set", "-n", "my-secret", "-g", "--exclude-special")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPotatoes))
	})

	It("generates a secret without numbers", func() {
		setupPostServer("my-secret", "potatoes", generateRequestJson(`{"exclude_number":true}`))

		session := runCommand("set", "-n", "my-secret", "-g", "--exclude-number")

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

		It("displays missing options message when neither generating or setting", func() {
			session := runCommand("set", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("One of the flags 'v' or 'g' must be specified"))
		})

		It("displays missing 'n' option as required parameter when only 'v' flag supplied", func() {
			session := runCommand("set", "-v", "potatoes")

			Eventually(session).Should(Exit(1))
			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
			}
		})

		It("displays missing 'n' option as required parameters", func() {
			session := runCommand("set")

			Eventually(session).Should(Exit(1))

			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
			}
		})

		It("displays unknown type message when given unknown type", func() {
			session := runCommand("set", "-n", "my-secret", "-v", "potatoes", "-t", "foobar")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("The request does not include a valid type. Please validate your input and retry your request."))
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

func setupPostServer(name string, value string, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(RESPONSE_JSON, value)),
		),
	)
}

func generateRequestJson(params string) string {
	return fmt.Sprintf(GENERATE_REQUEST_JSON, params)
}
