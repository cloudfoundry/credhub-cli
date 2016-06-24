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

var _ = Describe("Ca-Set", func() {
	It("puts a root CA", func() {
		var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "my-ca", "my-pub", "my-priv")
		setupPutCaServer("my-ca", "my-pub", "my-priv")

		session := runCommand("ca-set", "-n", "my-ca", "--public-string", "my-pub", "--private-string", "my-priv")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyCertificate))
	})

	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("ca-set", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("ca-set"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("public-string"))
			Expect(session.Err).To(Say("private-string"))
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

		It("displays the server provided error when an error is received", func() {
			server.AppendHandlers(
				RespondWith(http.StatusBadRequest, `{"error": "you fail."}`),
			)

			session := runCommand("ca-set", "-n", "my-ca", "--public-string", "my-pub", "--private-string", "my-priv")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupPutCaServer(name string, pub string, priv string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/ca/%s", name)),
			VerifyJSON(fmt.Sprintf(CA_REQUEST_JSON, pub, priv)),
			RespondWith(http.StatusOK, fmt.Sprintf(CA_RESPONSE_JSON, pub, priv)),
		),
	)
}
