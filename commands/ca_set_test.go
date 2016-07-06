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
			var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-ca", "my-pub", "my-priv")
			setupPutCaServer("root", "my-ca", "my-pub", "my-priv")

			session := runCommand("ca-set", "-n", "my-ca", "-t", "root", "--public-string", "my-pub", "--private-string", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("sets the type as root if no type is given", func() {
			var responseMyCertificate = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-ca", "my-pub", "my-priv")
			setupPutCaServer("root", "my-ca", "my-pub", "my-priv")

			session := runCommand("ca-set", "-n", "my-ca", "--public-string", "my-pub", "--private-string", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificate))
		})

		It("puts a secret using explicit certificate type and values read from files", func() {
			setupPutCaServer("root", "my-secret", "my-pub", "my-priv")
			tempDir := createTempDir("certFilesForTesting")
			publicFilename := createSecretFile(tempDir, "public.txt", "my-pub")
			privateFilename := createSecretFile(tempDir, "private.txt", "my-priv")

			session := runCommand("ca-set",
				"-n", "my-secret",
				"-t", "root",
				"--public", publicFilename,
				"--private", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseMyCertificateAuthority))
		})

		It("fails to put a CA when failing to read from a file", func() {
			testCaSetFileFailure("", "private.txt")
			testCaSetFileFailure("public.txt", "")
		})

		It("fails to put a CA when a specified cert string duplicates the contents of a file", func() {
			testSetCaFileDuplicationFailure("--public-string", "my-pub")
			testSetCaFileDuplicationFailure("--private-string", "my-priv")
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

func setupPutCaServer(caType, name, pub, priv string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", fmt.Sprintf("/api/v1/ca/%s", name)),
			VerifyJSON(fmt.Sprintf(CA_REQUEST_JSON, caType, pub, priv)),
			RespondWith(http.StatusOK, fmt.Sprintf(CA_RESPONSE_JSON, caType, pub, priv)),
		),
	)
}

func testCaSetFileFailure(publicFilename, privateFilename string) {
	tempDir := createTempDir("certFilesForTesting")
	if publicFilename != "" {
		publicFilename = createSecretFile(tempDir, publicFilename, "my-pub")
	} else {
		publicFilename = "dud"
	}
	if privateFilename != "" {
		privateFilename = createSecretFile(tempDir, privateFilename, "my-priv")
	} else {
		privateFilename = "dud"
	}

	session := runCommand("ca-set", "-n", "my-ca",
		"--public", publicFilename, "--private", privateFilename)

	os.RemoveAll(tempDir)
	Eventually(session).Should(Exit(1))
	Eventually(session.Err).Should(Say("A referenced file could not be opened. Please validate the provided filenames and permissions, then retry your request."))
}

func testSetCaFileDuplicationFailure(option, optionValue string) {
	setupPutCaServer("root", "my-secret", "my-pub", "my-priv")
	tempDir := createTempDir("certFilesForTesting")
	publicFilename := createSecretFile(tempDir, "public.txt", "my-pub")
	privateFilename := createSecretFile(tempDir, "private.txt", "my-priv")

	session := runCommand("ca-set",
		"-n", "my-secret",
		"-t", "root",
		"--public", publicFilename,
		"--private", privateFilename,
		option, optionValue)

	os.RemoveAll(tempDir)
	Eventually(session).Should(Exit(1))
	Eventually(session.Err).Should(Say("The combination of parameters in the request is not allowed. Please validate your input and retry your request."))
}
