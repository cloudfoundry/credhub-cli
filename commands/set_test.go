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

		It("puts a secret using explicit password type", func() {
			setupPutValueServer("my-password", "password", "potatoes")

			session := runCommand("set", "-n", "my-password", "-v", "potatoes", "-t", "password")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyPasswordPotatoes))
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
		It("displays help", func() {
			session := runCommand("set", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("set"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("secret"))
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
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(fmt.Sprintf(SECRET_STRING_REQUEST_JSON, secretType, value)),
			RespondWith(http.StatusOK, fmt.Sprintf(SECRET_STRING_RESPONSE_JSON, secretType, value)),
		),
	)
}

func setupPutCertificateServer(name string, ca string, cert string, priv string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(fmt.Sprintf(SECRET_CERTIFICATE_REQUEST_JSON, ca, cert, priv)),
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
