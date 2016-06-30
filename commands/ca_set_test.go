package commands_test

import (
	"net/http"

	"fmt"

	"runtime"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Ca-Set", func() {
	Describe("setting certificate authorities", func() {
		It("puts a root CA", func() {
			var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "my-ca", "my-pub", "my-priv")
			setupPutCaServer("my-ca", "my-pub", "my-priv")

			session := runCommand("ca-set", "-n", "my-ca", "--public-string", "my-pub", "--private-string", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("puts a secret using explicit certificate type and values read from files", func() {
			setupPutCaServer("my-secret", "my-pub", "my-priv")
			tempDir := createTempDir("certFilesForTesting")
			publicFilename := createSecretFile(tempDir, "public.txt", "my-pub")
			privateFilename := createSecretFile(tempDir, "private.txt", "my-priv")

			session := runCommand("ca-set",
				"-n", "my-secret",
				"--public", publicFilename,
				"--private", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificateAuthority))
		})
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
			session := runCommand("ca-set")

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
