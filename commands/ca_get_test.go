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
	ItBehavesLikeHelp("ca-get", "cg", func(session *Session) {
		Expect(session.Err).To(Say("ca-get command options"))
	})

	Describe("getting certificate authorities", func() {
		It("gets a root CA", func() {
			var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-ca-name", "my-cert", "my-priv")
			setupGetCaServer("root", "my-ca-name", "my-cert", "my-priv")

			session := runCommand("ca-get", "-n", "my-ca-name")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("can output a root CA as JSON", func() {
			setupGetCaServer("root", "my-ca-name", "my-cert", "my-priv")

			session := runCommand("ca-get", "-n", "my-ca-name", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(`{
				"type": "root",
				"updated_at": "` + TIMESTAMP + `",
				"certificate": "my-cert",
				"private_key": "my-priv"
			}`))
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

func setupGetCaServer(caType, name, cert, priv string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("GET", "/api/v1/ca", "name="+name+"&current=true"),
			RespondWith(http.StatusOK, fmt.Sprintf(CA_ARRAY_RESPONSE_JSON, caType, cert, priv)),
		),
	)
}
