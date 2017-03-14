package commands_test

import (
	"net/http"

	"fmt"

	"runtime"

	"os"

	"strings"

	"github.com/cloudfoundry-incubator/credhub-cli/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Set", func() {
	Describe("error handling", func() {
		It("displays an error when an invalid type is specified", func() {
			server.AppendHandlers(
				RespondWith(http.StatusBadRequest, `{"error": "we should not reach this point"}`),
			)

			session := runCommand("set", "-t", "potato", "-n", "foo")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("The request does not include a valid type. Valid values include 'value', 'password', 'certificate', 'ssh', 'rsa', and 'json'."))
		})
	})

	Describe("setting string secrets", func() {
		It("puts a secret using explicit value type", func() {
			setupPutValueServer("my-value", "value", "potatoes")

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyValuePotatoes))
		})
	})

	Describe("setting json secrets", func() {
		It("puts a secret using explicit json type", func() {
			jsonValue := `{"foo":"bar","nested":{"a":1},"an":["array"]}`
			setupPutJsonServer("json-secret", jsonValue)

			session := runCommand("set", "-n", "json-secret", "-v", jsonValue, "-t", "json")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyJson))
		})
	})

	Describe("setting SSH secrets", func() {
		It("puts a secret using explicit ssh type", func() {
			setupPutRsaSshServer("foo-ssh-key", "ssh", "some-public-key", "some-private-key", true)

			session := runCommand("set", "-n", "foo-ssh-key", "-U", "some-public-key", "-P", "some-private-key", "-t", "ssh")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMySSHFoo))
		})

		It("puts a secret using values read from files", func() {
			setupPutRsaSshServer("foo-ssh-key", "ssh", "some-public-key", "some-private-key", true)

			tempDir := createTempDir("sshFilesForTesting")
			publicFileName := createSecretFile(tempDir, "rsa.pub", "some-public-key")
			privateFilename := createSecretFile(tempDir, "rsa.key", "some-private-key")

			session := runCommand("set", "-n", "foo-ssh-key",
				"-t", "ssh",
				"-u", publicFileName,
				"-p", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMySSHFoo))
		})

		It("puts a secret specifying no-overwrite", func() {
			setupPutRsaSshServer("foo-ssh-key", "ssh", "some-public-key", "some-private-key", false)

			session := runCommand("set", "-n", "foo-ssh-key", "-t", "ssh", "-U", "some-public-key", "-P", "some-private-key", "--no-overwrite")

			Eventually(session).Should(Exit(0))
		})
	})

	Describe("setting RSA secrets", func() {
		It("puts a secret using explicit rsa type", func() {
			setupPutRsaSshServer("foo-rsa-key", "rsa", "some-public-key", "some-private-key", true)

			session := runCommand("set", "-n", "foo-rsa-key", "-U", "some-public-key", "-P", "some-private-key", "-t", "rsa")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyRSAFoo))
		})

		It("puts a secret using values read from files", func() {
			setupPutRsaSshServer("foo-rsa-key", "rsa", "some-public-key", "some-private-key", true)

			tempDir := createTempDir("rsaFilesForTesting")
			publicFileName := createSecretFile(tempDir, "rsa.pub", "some-public-key")
			privateFilename := createSecretFile(tempDir, "rsa.key", "some-private-key")

			session := runCommand("set", "-n", "foo-rsa-key",
				"-t", "rsa",
				"-u", publicFileName,
				"-p", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyRSAFoo))
		})

		It("puts a secret specifying no-overwrite", func() {
			setupPutRsaSshServer("foo-rsa-key", "rsa", "some-public-key", "some-private-key", false)

			session := runCommand("set", "-n", "foo-rsa-key", "-t", "rsa", "-U", "some-public-key", "-P", "some-private-key", "--no-overwrite")

			Eventually(session).Should(Exit(0))
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

			response := fmt.Sprintf(STRING_SECRET_RESPONSE_YAML, "password", "my-password", "potatoes potatoes")

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

func setupPutRsaSshServer(name, keyType, publicKey, privateKey string, overwrite bool) {
	var jsonRequest string
	jsonRequest = fmt.Sprintf(RSA_SSH_SECRET_REQUEST_JSON, keyType, name, publicKey, privateKey, overwrite)
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(jsonRequest),
			RespondWith(http.StatusOK, fmt.Sprintf(RSA_SSH_SECRET_RESPONSE_JSON, keyType, name, publicKey, privateKey)),
		),
	)
}

func setupPutValueServer(name, secretType, value string) {
	setupOverwritePutValueServer(name, secretType, value, true)
}

func setupOverwritePutValueServer(name, secretType, value string, overwrite bool) {
	var jsonRequest string
	jsonRequest = fmt.Sprintf(STRING_SECRET_OVERWRITE_REQUEST_JSON, secretType, name, value, overwrite)
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(jsonRequest),
			RespondWith(http.StatusOK, fmt.Sprintf(STRING_SECRET_RESPONSE_JSON, secretType, name, value)),
		),
	)
}

func setupPutJsonServer(name, value string) {
	setupOverwritePutJsonServer(name, value, true)
}

func setupOverwritePutJsonServer(name, value string, overwrite bool) {
	var jsonRequest string
	jsonRequest = fmt.Sprintf(JSON_SECRET_OVERWRITE_REQUEST_JSON, name, value, overwrite)
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(jsonRequest),
			RespondWith(http.StatusOK, fmt.Sprintf(JSON_SECRET_RESPONSE_JSON, name, value)),
		),
	)
}

func setupPutCertificateServer(name, ca, cert, priv string) {
	setupOverwritePutCertificateServer(name, ca, cert, priv, true)
}

func setupOverwritePutCertificateServer(name, ca, cert, priv string, overwrite bool) {
	var jsonRequest string
	jsonRequest = fmt.Sprintf(CERTIFICATE_SECRET_REQUEST_JSON, name, ca, cert, priv, overwrite)
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(jsonRequest),
			RespondWith(http.StatusOK, fmt.Sprintf(CERTIFICATE_SECRET_RESPONSE_JSON, name, ca, cert, priv)),
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
