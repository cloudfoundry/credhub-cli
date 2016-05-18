package commands_test

import (
	"net/http"

	"fmt"

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
		Expect(session.Err).To(Say("get"))
		Expect(session.Err).To(Say("name"))
	})

	It("gets a secret", func() {
		responseJson := `{"value":"potatoes"}`
		responseTable := fmt.Sprintf(`Name:	my-secret\nValue:	potatoes`)

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
