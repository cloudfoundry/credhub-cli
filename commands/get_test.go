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

	It("gets a secret", func() {
		responseJson := `{"type":"value","value":"potatoes"}`
		responseTable := fmt.Sprintf(`Type:		value\nName:		my-secret\nValue:	potatoes`)

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data/my-secret"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseTable))
	})
})
