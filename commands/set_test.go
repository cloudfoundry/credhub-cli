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
	It("puts a secret", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)
		requestJson := `{"value":"potatoes"}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("PUT", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-v", "potatoes")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	It("generates a secret without parameters", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)
		requestJson := `{"parameters":{}}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("POST", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-g")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	It("generates a secret with length", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)
		requestJson := `{"parameters":{"length":42}}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("POST", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-g", "-l", "42")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	It("generates a secret without upper case", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)
		requestJson := `{"parameters":{"exclude_upper":true}}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("POST", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-g", "--exclude-upper")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	It("generates a secret without lower case", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)
		requestJson := `{"parameters":{"exclude_lower":true}}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("POST", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-g", "--exclude-lower")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	It("generates a secret without special characters", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)
		requestJson := `{"parameters":{"exclude_special":true}}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("POST", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-g", "--exclude-special")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	It("generates a secret without numbers", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)
		requestJson := `{"parameters":{"exclude_number":true}}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("POST", "/api/v1/data/my-secret"),
				VerifyJSON(requestJson),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("set", "-n", "my-secret", "-g", "--exclude-number")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
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
	})
})
