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
	"github.com/pivotal-cf/credhub-cli/commands"
)

var _ = Describe("Ca-Set", func() {
	Describe("setting certificate authorities", func() {
		It("puts a root CA", func() {
			var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-ca", "my-cert", "my-priv")
			setupPutCaServer("root", "my-ca", "my-cert", "my-priv")

			session := runCommand("ca-set", "-n", "my-ca", "-t", "root", "--certificate-string", "my-cert", "--private-string", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("puts a root CA", func() {
			setupPutCaServer("root", "my-ca", "my-cert", "my-priv")

			session := runCommand(
				"ca-set",
				"-n", "my-ca",
				"-t", "root",
				"--certificate-string", "my-cert",
				"--private-string", "my-priv",
				"--output-json",
			)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(`{
				"type": "root",
				"updated_at": "` + TIMESTAMP + `",
				"certificate": "my-cert",
				"private_key": "my-priv"
			}`))
		})

		It("sets the type as root if no type is given", func() {
			var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-ca", "my-cert", "my-priv")
			setupPutCaServer("root", "my-ca", "my-cert", "my-priv")

			session := runCommand("ca-set", "-n", "my-ca", "--certificate-string", "my-cert", "--private-string", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("puts a secret using explicit certificate type and values read from files", func() {
			setupPutCaServer("root", "my-secret", "my-cert", "my-priv")
			tempDir := createTempDir("certFilesForTesting")
			certificateFilename := createSecretFile(tempDir, "certificate.txt", "my-cert")
			privateFilename := createSecretFile(tempDir, "private.txt", "my-priv")

			session := runCommand("ca-set",
				"-n", "my-secret",
				"-t", "root",
				"--certificate", certificateFilename,
				"--private", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificateAuthority))
		})

		It("fails to put a CA when failing to read from a file", func() {
			testCaSetFileFailure("", "private.txt")
			testCaSetFileFailure("certificate.txt", "")
		})

		It("fails to put a CA when a specified cert string duplicates the contents of a file", func() {
			testSetCaFileDuplicationFailure("--certificate-string", "my-cert")
			testSetCaFileDuplicationFailure("--private-string", "my-priv")
		})
	})

	Describe("Help", func() {
		It("has short flags", func() {
			Expect(commands.CaSetCommand{}).To(SatisfyAll(
				commands.HaveFlag("name", "n"),
				commands.HaveFlag("type", "t"),
				commands.HaveFlag("certificate", "c"),
				commands.HaveFlag("private", "p"),
				commands.HaveFlag("certificate-string", "C"),
				commands.HaveFlag("private-string", "P"),
			))
		})

		ItBehavesLikeHelp("ca-set", "cs", func(session *Session) {
			Expect(session.Err).To(Say("ca-set"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("certificate-string"))
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

			session := runCommand("ca-set", "-n", "my-ca", "--certificate-string", "my-cert", "--private-string", "my-priv")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupPutCaServer(caType, name, certificate, priv string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/ca/%s", name)),
			VerifyJSON(fmt.Sprintf(CA_REQUEST_JSON, caType, certificate, priv)),
			RespondWith(http.StatusOK, fmt.Sprintf(CA_SET_RESPONSE_JSON, caType, certificate, priv)),
		),
	)
}

func testCaSetFileFailure(certificateFilename, privateFilename string) {
	tempDir := createTempDir("certFilesForTesting")
	if certificateFilename != "" {
		certificateFilename = createSecretFile(tempDir, certificateFilename, "my-cert")
	} else {
		certificateFilename = "dud"
	}
	if privateFilename != "" {
		privateFilename = createSecretFile(tempDir, privateFilename, "my-priv")
	} else {
		privateFilename = "dud"
	}

	session := runCommand("ca-set", "-n", "my-ca",
		"--certificate", certificateFilename, "--private", privateFilename)

	os.RemoveAll(tempDir)
	Eventually(session).Should(Exit(1))
	Eventually(session.Err).Should(Say("A referenced file could not be opened. Please validate the provided filenames and permissions, then retry your request."))
}

func testSetCaFileDuplicationFailure(option, optionValue string) {
	setupPutCaServer("root", "my-secret", "my-cert", "my-priv")
	tempDir := createTempDir("certFilesForTesting")
	certificateFilename := createSecretFile(tempDir, "certificate.txt", "my-cert")
	privateFilename := createSecretFile(tempDir, "private.txt", "my-priv")

	session := runCommand("ca-set",
		"-n", "my-secret",
		"-t", "root",
		"--certificate", certificateFilename,
		"--private", privateFilename,
		option, optionValue)

	os.RemoveAll(tempDir)
	Eventually(session).Should(Exit(1))
	Eventually(session.Err).Should(Say("The combination of parameters in the request is not allowed. Please validate your input and retry your request."))
}
