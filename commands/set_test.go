package commands_test

import (
	"net/http"

	"fmt"

	"runtime"

	"os"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/credhub-cli/commands"
)

var _ = Describe("Set", func() {
	Describe("setting string secrets", func() {
		It("puts a secret using explicit value type", func() {
			setupPutValueServer("my-value", "value", "potatoes")

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyValuePotatoes))
		})
	})

	Describe("setting password secrets", func() {
		It("puts a secret using default type", func() {
			setupPutValueServer("my-password", "password", "potatoes")

			session := runCommand("set", "-n", "my-password", "-v", "potatoes")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyPasswordPotatoes))
		})

		It("puts a secret specifying no-overwrite", func() {
			setupOverwritePutValueServer("my-password", "password", "potatoes", false)

			session := runCommand("set", "-n", "my-password", "-v", "potatoes", "--no-overwrite")

			Eventually(session).Should(Exit(0))
		})

		It("puts a secret using explicit password type", func() {
			setupPutValueServer("my-password", "password", "potatoes")

			session := runCommand("set", "-n", "my-password", "-v", "potatoes", "-t", "password")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyPasswordPotatoes))
		})

		It("prompts for value if value is not provided", func() {
			setupPutValueServer("my-password", "password", "potatoes")

			session := runCommandWithStdin(strings.NewReader("potatoes\n"), "set", "-n", "my-password", "-t", "password")

			Eventually(session.Out).Should(Say("value:"))
			Eventually(session.Wait("10s").Out).Should(Say(responseMyPasswordPotatoes))
			Eventually(session).Should(Exit(0))
		})

		It("can set password that contains spaces interactively", func() {
			setupPutValueServer("my-password", "password", "potatoes potatoes")

			session := runCommandWithStdin(strings.NewReader("potatoes potatoes\n"), "set", "-n", "my-password", "-t", "password")

			response := fmt.Sprintf(SECRET_STRING_RESPONSE_TABLE, "password", "my-password", "potatoes potatoes")

			Eventually(session.Out).Should(Say("value:"))
			Eventually(session.Wait("10s").Out).Should(Say(response))
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("setting certificate secrets", func() {
		It("puts a secret using explicit certificate type and string values", func() {
			setupPutCertificateServer("my-secret", "my-ca", "my-cert", "my-priv")

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root-string", "my-ca",
				"--certificate-string", "my-cert", "--private-string", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("puts a secret using explicit certificate type and string values with no-overwrite", func() {
			setupOverwritePutCertificateServer("my-secret", "my-ca", "my-cert", "my-priv", false)

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root-string", "my-ca",
				"--certificate-string", "my-cert", "--private-string", "my-priv", "--no-overwrite")

			Eventually(session).Should(Exit(0))
		})

		It("puts a secret using explicit certificate type and values read from files", func() {
			setupPutCertificateServer("my-secret", "my-ca", "my-cert", "my-priv")
			tempDir := createTempDir("certFilesForTesting")
			caFilename := createSecretFile(tempDir, "ca.txt", "my-ca")
			certificateFilename := createSecretFile(tempDir, "certificate.txt", "my-cert")
			privateFilename := createSecretFile(tempDir, "private.txt", "my-priv")

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", caFilename,
				"--certificate", certificateFilename, "--private", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("fails to put a secret when failing to read from a file", func() {
			testSetFileFailure("", "certificate.txt", "private.txt")
			testSetFileFailure("ca.txt", "", "private.txt")
			testSetFileFailure("ca.txt", "certificate.txt", "")
		})

		It("fails to put a secret when a specified cert string duplicates the contents of a file", func() {
			testSetCertFileDuplicationFailure("--root-string", "my-ca")
			testSetCertFileDuplicationFailure("--certificate-string", "my-cert")
			testSetCertFileDuplicationFailure("--private-string", "my-priv")
		})
	})

	Describe("Help", func() {
		It("short flags", func() {
			Expect(commands.SetCommand{}).To(SatisfyAll(
				commands.HaveFlag("name", "n"),
				commands.HaveFlag("type", "t"),
				commands.HaveFlag("value", "v"),
				commands.HaveFlag("no-overwrite", "O"),
				commands.HaveFlag("root", "r"),
				commands.HaveFlag("certificate", "c"),
				commands.HaveFlag("private", "p"),
				commands.HaveFlag("root-string", "R"),
				commands.HaveFlag("certificate-string", "C"),
				commands.HaveFlag("private-string", "P"),
			))
		})

		ItBehavesLikeHelp("set", "s", func(session *Session) {
			Expect(session.Err).To(Say("set"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("credential"))
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

			session := runCommand("set", "-n", "my-value", "-v", "tomatoes")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupPutValueServer(name string, secretType string, value string) {
	setupOverwritePutValueServer(name, secretType, value, true)
}

func setupOverwritePutValueServer(name string, secretType string, value string, overwrite bool) {
	var jsonRequest string
	jsonRequest = fmt.Sprintf(SECRET_STRING_OVERWRITE_REQUEST_JSON, secretType, value, overwrite)
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(jsonRequest),
			RespondWith(http.StatusOK, fmt.Sprintf(SECRET_STRING_RESPONSE_JSON, secretType, value)),
		),
	)
}

func setupPutCertificateServer(name string, ca string, cert string, priv string) {
	setupOverwritePutCertificateServer(name, ca, cert, priv, true)
}

func setupOverwritePutCertificateServer(name string, ca string, cert string, priv string, overwrite bool) {
	var jsonRequest string
	jsonRequest = fmt.Sprintf(SECRET_CERTIFICATE_REQUEST_JSON, ca, cert, priv, overwrite)
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(jsonRequest),
			RespondWith(http.StatusOK, fmt.Sprintf(SECRET_CERTIFICATE_RESPONSE_JSON, ca, cert, priv)),
		),
	)
}

func testSetFileFailure(caFilename, certificateFilename, privateFilename string) {
	tempDir := createTempDir("certFilesForTesting")
	if caFilename != "" {
		caFilename = createSecretFile(tempDir, caFilename, "my-ca")
	} else {
		caFilename = "dud"
	}
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

	session := runCommand("set", "-n", "my-secret",
		"-t", "certificate", "--root", caFilename,
		"--certificate", certificateFilename, "--private", privateFilename)

	os.RemoveAll(tempDir)
	Eventually(session).Should(Exit(1))
	Eventually(session.Err).Should(Say("A referenced file could not be opened. Please validate the provided filenames and permissions, then retry your request."))
}

func testSetCertFileDuplicationFailure(option, optionValue string) {
	setupPutCertificateServer("my-secret", "my-ca", "my-cert", "my-priv")
	tempDir := createTempDir("certFilesForTesting")
	caFilename := createSecretFile(tempDir, "ca.txt", "my-ca")
	certificateFilename := createSecretFile(tempDir, "certificate.txt", "my-cert")
	privateFilename := createSecretFile(tempDir, "private.txt", "my-priv")

	session := runCommand("set", "-n", "my-secret", "-t", "certificate", "--root", caFilename,
		"--certificate", certificateFilename, "--private", privateFilename, option, optionValue)

	os.RemoveAll(tempDir)
	Eventually(session).Should(Exit(1))
	Eventually(session.Err).Should(Say("The combination of parameters in the request is not allowed. Please validate your input and retry your request."))
}
