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
	It("displays help", func() {
		session := runCommand("set", "-h")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("set"))
		Expect(session.Err).To(Say("name"))
		Expect(session.Err).To(Say("secret"))
	})

	Describe("Flags", func() {
		It("displays missing 's' option as required parameter", func() {
			session := runCommand("set", "-n", "my-secret")

			Eventually(session).Should(Exit(1))
			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/s, /secret' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-s, --secret' was not specified"))
			}
		})

		It("displays missing 'n' option as required parameter", func() {
			session := runCommand("set", "-s", "potatoes")

			Eventually(session).Should(Exit(1))
			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
			}
		})

		It("displays missing 'n' option and 's' option as required parameters", func() {
			session := runCommand("set")

			Eventually(session).Should(Exit(1))

			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flags `/n, /name' and `/s, /secret' were not specified"))
			} else {
				Expect(session.Err).To(Say("the required flags `-n, --name' and `-s, --secret' were not specified"))
			}
		})
	})

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

		session := runCommand("set", "-n", "my-secret", "-s", "potatoes")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})
})
