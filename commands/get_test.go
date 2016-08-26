package commands_test

import (
	"net/http"

	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Get", func() {
	It("displays help", func() {
		session := runCommand("get", "-h")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("Usage"))
		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("cm-cli.exe \\[OPTIONS\\] get \\[get-OPTIONS\\]"))
		} else {
			Expect(session.Err).To(Say("cm-cli \\[OPTIONS\\] get \\[get-OPTIONS\\]"))
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
		responseJson := `{"type":"value","value":"potatoes"}`
		responseTable := `Type:		value\nName:		my-value\nValue:		potatoes`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data/my-value"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-value")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})

	It("gets a password secret", func() {
		responseJson := `{"type":"password","value":"potatoes"}`
		responseTable := `Type:		password\nName:		my-password\nValue:		potatoes`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data/my-password"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-password")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})
})
