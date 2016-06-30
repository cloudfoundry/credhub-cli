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

var _ = Describe("Ca-Get", func() {
	Describe("getting certificate authorities", func() {
		It("gets a root CA", func() {
			var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-ca-name", "my-pub", "my-priv")
			setupGetCaServer("root", "my-ca-name", "my-pub", "my-priv")

			session := runCommand("ca-get", "-n", "my-ca-name")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("displays the server provided error if it cannot get ca by name", func() {
			server.AppendHandlers(
				RespondWith(http.StatusBadRequest, `{"error": "you fail."}`),
			)
			session := runCommand("ca-get", "-n", "my-ca-name")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupGetCaServer(caType, name, pub, priv string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("GET", fmt.Sprintf("/api/v1/ca/%s", name)),
			RespondWith(http.StatusOK, fmt.Sprintf(CA_RESPONSE_JSON, caType, pub, priv)),
		),
	)
}
