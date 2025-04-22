package commands_test

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"code.cloudfoundry.org/credhub-cli/commands"
	"code.cloudfoundry.org/credhub-cli/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Set", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("set", "-n", "test-credential", "-t", "password", "-w", "value")
	ItRequiresAnAPIToBeSet("set", "-n", "test-credential", "-t", "password", "-w", "value")
	testAutoLogin := []TestAutoLogin{
		{
			method:              "PUT",
			responseFixtureFile: "set_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v1/data",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogin, "set", "-n", "test-credential", "-t", "password", "-w", "test-value")

	It("returns an error when not specifying type", func() {
		session := runCommand("set", "-n", "my-password", "-w", "potatoes")

		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("A type must be specified when setting a credential. Valid types include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'."))
	})

	It("returns an error when metadata json is invalid", func() {
		setupSetServerWithMetadata("my-value", "value", `"potatoes"`, `"not-valid-json"`)

		session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value", "--metadata", "not-valid-json")

		Eventually(session).Should(Exit(1))
		Expect(string(session.Err.Contents())).To(ContainSubstring("The argument for --metadata is not a valid json object. Please update and retry your request."))
	})

	It("errors when server does not support metadata", func() {
		setCachedServerVersion("2.5.0")

		session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

		Eventually(session).Should(Exit(1))
		Expect(string(session.Err.Contents())).To(ContainSubstring("The --metadata flag is not supported for this version of the credhub server (requires >= 2.6.x). Please remove the flag and retry your request."))
	})

	Describe("setting value secrets", func() {
		It("puts a secret", func() {
			setupSetServer("my-value", "value", `"potatoes"`)

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring("name: my-value"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("type: value"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("value: <redacted>"))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("metadata:"))
		})

		It("puts a secret with metadata", func() {
			setupSetServerWithMetadata("my-value", "value", `"potatoes"`, `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(0), string(session.Err.Contents()))
			metadataOutput := `
metadata:
    array:
        - metadata
    some:
        example: metadata`
			Expect(string(session.Out.Contents())).To(ContainSubstring(metadataOutput))
		})

		It("escapes special characters in the value", func() {
			setupSetServer("my-character-test", "value", `"{\"password\":\"some-still-bad-password\"}"`)

			session := runCommand("set", "-t", "value", "-n", "my-character-test", "-v", `{"password":"some-still-bad-password"}`)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-character-test"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: value"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`value: <redacted>`))
		})

		It("puts a secret and returns in json format", func() {
			setupSetServer("my-value", "value", `"potatoes"`)

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(fmt.Sprintf(redactedResponseJSON, "value", "my-value", "null")))
		})

		It("accepts case-insensitive type", func() {
			setupSetServer("my-value", "value", `"potatoes"`)

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "VALUE", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(fmt.Sprintf(redactedResponseJSON, "value", "my-value", "null")))
		})
	})

	Describe("setting json secrets", func() {
		It("puts a secret", func() {
			jsonValue := `{"foo":"bar","nested":{"a":1},"an":["array"]}`
			setupSetServer("json-secret", "json", jsonValue)

			session := runCommand("set", "-n", "json-secret", "-v", jsonValue, "-t", "json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring("name: json-secret"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("type: json"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("value: <redacted>"))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("metadata:"))
		})

		It("puts a secret with metadata", func() {
			jsonValue := `{"foo":"bar","nested":{"a":1},"an":["array"]}`

			setupSetServerWithMetadata("json-secret", "json", jsonValue, `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			session := runCommand("set", "-n", "json-secret", "-v", jsonValue, "-t", "json", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(0))
			metadataOutput := `
metadata:
    array:
        - metadata
    some:
        example: metadata`
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(metadataOutput))
		})

		It("escapes special characters in the json", func() {
			setupSetServer("my-character-test", "json", `{"foo":"b\"ar"}`)

			session := runCommand("set", "-t", "json", "-n", "my-character-test", "-v", `{"foo":"b\"ar"}`)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-character-test"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: json"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret and returns in json format", func() {
			jsonValue := `{"foo":"bar","nested":{"a":1},"an":["array"]}`
			setupSetServer("json-secret", "json", jsonValue)

			session := runCommand("set", "-n", "json-secret", "-v", jsonValue, "-t", "json", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(fmt.Sprintf(redactedResponseJSON, "json", "json-secret", "null")))
		})

		It("accepts case-insensitive type", func() {
			jsonValue := `{"foo":"bar","nested":{"a":1},"an":["array"]}`
			setupSetServer("json-secret", "json", jsonValue)

			session := runCommand("set", "-n", "json-secret", "-v", jsonValue, "-t", "JSON")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: json-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: json"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})
	})

	Describe("setting SSH secrets", func() {
		It("puts a secret", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-ssh-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "ssh")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring("name: foo-ssh-key"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("type: ssh"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("value: <redacted>"))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("metadata:"))
		})

		It("puts a secret with metadata", func() {
			setupSetServerWithMetadata("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			session := runCommand("set", "-n", "foo-ssh-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "ssh", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(0))
			metadataOutput := `
metadata:
    array:
        - metadata
    some:
        example: metadata`
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(metadataOutput))
		})

		It("puts a secret using values read from files", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			tempDir := test.CreateTempDir("sshFilesForTesting")
			publicFileName := test.CreateCredentialFile(tempDir, "rsa.pub", "some-public-key")
			privateFilename := test.CreateCredentialFile(tempDir, "rsa.key", "some-private-key")

			session := runCommand("set", "-n", "foo-ssh-key",
				"-t", "ssh",
				"-u", publicFileName,
				"-p", privateFilename)

			Expect(os.RemoveAll(tempDir)).To(Succeed())
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-ssh-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: ssh"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret and returns in json format", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-ssh-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "ssh", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(fmt.Sprintf(redactedResponseJSON, "ssh", "foo-ssh-key", "null")))
		})

		It("accepts case-insensitive type", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-ssh-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "SSH")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-ssh-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: ssh"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("handles newline characters", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some\npublic\nkey","private_key":"some\nprivate\nkey"}`)
			session := runCommand("set", "-n", "foo-ssh-key", "-u", `some\npublic\nkey`, "-p", `some\nprivate\nkey`, "-t", "ssh", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(fmt.Sprintf(redactedResponseJSON, "ssh", "foo-ssh-key", "null")))
		})
	})

	Describe("setting RSA secrets", func() {
		It("puts a secret ", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`)
			session := runCommand("set", "-n", "foo-rsa-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "rsa")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring("name: foo-rsa-key"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("type: rsa"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("value: <redacted>"))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("metadata:"))
		})

		It("puts a secret with metadata", func() {
			setupSetServerWithMetadata("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			session := runCommand("set", "-n", "foo-rsa-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "rsa", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(0), string(session.Err.Contents()))
			metadataOutput := `
metadata:
    array:
        - metadata
    some:
        example: metadata`
			Expect(string(session.Out.Contents())).To(ContainSubstring(metadataOutput))
		})

		It("puts a secret using values read from files", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			tempDir := test.CreateTempDir("rsaFilesForTesting")
			publicFileName := test.CreateCredentialFile(tempDir, "rsa.pub", "some-public-key")
			privateFilename := test.CreateCredentialFile(tempDir, "rsa.key", "some-private-key")

			session := runCommand("set", "-n", "foo-rsa-key",
				"-t", "rsa",
				"-u", publicFileName,
				"-p", privateFilename)

			Expect(os.RemoveAll(tempDir)).To(Succeed())
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-rsa-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: rsa"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret and returns in json format", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-rsa-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "rsa", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"name": "foo-rsa-key"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"type": "rsa"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"value": "<redacted>"`))

		})

		It("accepts case-insensitive type", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-rsa-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "RSA")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-rsa-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: rsa"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("handles newline characters", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some\npublic\nkey","private_key":"some\nprivate\nkey"}`)

			session := runCommand("set", "-n", "foo-rsa-key", "-u", `some\npublic\nkey`, "-p", `some\nprivate\nkey`, "-t", "rsa", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).Should(MatchJSON(fmt.Sprintf(redactedResponseJSON, "rsa", "foo-rsa-key", "null")))
		})
	})

	Describe("setting password secrets", func() {
		It("puts a secret and returns in yaml format", func() {
			setupSetServer("my-password", "password", `"potatoes"`)

			session := runCommand("set", "-n", "my-password", "-w", "potatoes", "-t", "password")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("metadata:"))
		})

		It("puts a secret with metadata", func() {
			setupSetServerWithMetadata("my-password", "password", `"potatoes"`, `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			session := runCommand("set", "-n", "my-password", "-w", "potatoes", "-t", "password", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(0), string(session.Err.Contents()))
			metadataOutput := `
metadata:
    array:
        - metadata
    some:
        example: metadata`
			Expect(string(session.Out.Contents())).To(ContainSubstring(metadataOutput))
		})

		It("prompts for value if value is not provided", func() {
			setupSetServer("my-password", "password", `"potatoes"`)

			session := runCommandWithStdin(strings.NewReader("potatoes\n"), "set", "-n", "my-password", "-t", "password")

			Eventually(string(session.Out.Contents())).Should(ContainSubstring("password: ********"))
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("can set password that contains spaces interactively", func() {
			setupSetServer("my-password", "password", `"potatoes potatoes"`)

			session := runCommandWithStdin(strings.NewReader("potatoes potatoes\n"), "set", "-n", "my-password", "-t", "password")

			Eventually(string(session.Out.Contents())).Should(ContainSubstring("password:"))
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("escapes special characters in the password", func() {
			setupSetServer("my-password", "password", `"{\"password\":\"some-still-bad-password\"}"`)

			session := runCommand("set", "-t", "password", "-n", "my-password", "-w", `{"password":"some-still-bad-password"}`)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`value: <redacted>`))
		})

		It("puts a secret and returns in json format", func() {
			setupSetServer("my-password", "password", `"potatoes"`)

			session := runCommand("set", "-n", "my-password", "-w", "potatoes", "-t", "password", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(fmt.Sprintf(redactedResponseJSON, "password", "my-password", "null")))
		})

		It("accepts case-insensitive type", func() {
			setupSetServer("my-password", "password", `"potatoes"`)

			session := runCommand("set", "-n", "my-password", "-w", "potatoes", "-t", "PASSWORD")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})
	})

	Describe("setting certificate secrets", func() {
		It("puts a secret and string values", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)
			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: certificate"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("metadata:"))
		})

		It("puts a secret with metadata", func() {
			setupSetServerWithMetadata("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`, `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv",
				"--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(0), string(session.Err.Contents()))
			metadataOutput := `
metadata:
    array:
        - metadata
    some:
        example: metadata`
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(metadataOutput))
		})

		It("puts a secret, string values, and certificate authority name", func() {
			setupSetServer("my-secret", "certificate", `{"ca": "", "ca_name":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--ca-name", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: certificate"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret and values read from files", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)
			tempDir := test.CreateTempDir("certFilesForTesting")
			caFilename := test.CreateCredentialFile(tempDir, "ca.txt", "my-ca")
			certificateFilename := test.CreateCredentialFile(tempDir, "certificate.txt", "my-cert")
			privateFilename := test.CreateCredentialFile(tempDir, "private.txt", "my-priv")

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", caFilename,
				"--certificate", certificateFilename, "--private", privateFilename)

			Expect(os.RemoveAll(tempDir)).To(Succeed())
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: certificate"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		if runtime.GOOS != "windows" {
			It("fails to put a secret when reading from unreadable file", func() {
				testSetFileFailure("unreadable.txt", "", "")
				testSetFileFailure("", "unreadable.txt", "")
				testSetFileFailure("", "", "unreadable.txt")
			})
		}

		It("puts a secret and string values in json format", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"name": "my-secret"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"type": "certificate"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"value": "<redacted>"`))
		})

		It("accepts case insensitive type", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)

			session := runCommand("set", "-n", "my-secret",
				"-t", "CERTIFICATE", "--root", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: certificate"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("handles newline characters", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my\nca","certificate":"my\ncert","private_key":"my\npriv"}`)
			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", `my\nca`,
				"--certificate", `my\ncert`, "--private", `my\npriv`, "--output-json")
			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).Should(MatchJSON(fmt.Sprintf(redactedResponseJSON, "certificate", "my-secret", "null")))
		})
	})

	Describe("setting User secrets", func() {
		It("puts a secret", func() {
			setupSetServer("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`)
			session := runCommand("set", "-n", "my-username-credential", "-z", "my-username", "-w", "test-password", "-t", "user")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-username-credential"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: user"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`version_created_at: "2016-01-01T12:00:00Z"`))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("metadata:"))
		})

		It("puts a secret with metadata", func() {
			setupSetServerWithMetadata("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`, `{"some":{"example":"metadata"}, "array":["metadata"]}`)
			session := runCommand("set", "-n", "my-username-credential", "-z", "my-username", "-w", "test-password", "-t", "user", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(0), string(session.Err.Contents()))
			metadataOutput := `
metadata:
    array:
        - metadata
    some:
        example: metadata`
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(metadataOutput))
		})

		It("should set password interactively for user", func() {
			setupSetServer("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`)
			session := runCommandWithStdin(strings.NewReader("test-password\n"), "set", "-n", "my-username-credential", "-t", "user", "--username", "my-username")

			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-username-credential"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: user"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
			Eventually(session).Should(Exit(0))
		})

		It("should set null username when it isn't provided", func() {
			setupSetServer("my-username-credential", "user", `{"username": "", "password": "test-password"}`)

			session := runCommandWithStdin(strings.NewReader("test-password\n"), "set", "-n", "my-username-credential", "-t", "user")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-username-credential"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: user"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret in json format", func() {
			setupSetServer("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`)

			session := runCommand("set", "-n", "my-username-credential", "-z", "my-username", "-w", "test-password", "-t", "user",
				"--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"name": "my-username-credential"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"type": "user"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"value": "<redacted>"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"version_created_at": "2016-01-01T12:00:00Z"`))
		})

		It("accepts case-insensitive type", func() {
			setupSetServer("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`)
			session := runCommand("set", "-n", "my-username-credential", "-z", "my-username", "-w", "test-password", "-t", "USER")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-username-credential"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: user"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`version_created_at: "2016-01-01T12:00:00Z"`))
		})
	})

	Describe("Help", func() {
		It("short flags", func() {
			Expect(commands.SetCommand{}).To(SatisfyAll(
				commands.HaveFlag("name", "n"),
				commands.HaveFlag("type", "t"),
				commands.HaveFlag("value", "v"),
				commands.HaveFlag("root", "r"),
				commands.HaveFlag("certificate", "c"),
				commands.HaveFlag("private", "p"),
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
				RespondWith(http.StatusBadRequest, `{"error": "test error"}`),
			)

			session := runCommand("set", "-n", "my-value", "-t", "value", "-v", "tomatoes")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("test error"))
		})
	})
})

const setCredentialRequestJSON = `{"type":"%s","name":"%s","value":%s}`
const setCredentialResponseJSON = `{"type":"%s","id":"` + uuid + `","name":"%s","value":%s,"version_created_at":"` + timestamp + `"}`
const setCredentialRequestJSONWithMetadata = `{"type":"%s","name":"%s","value":%s,"metadata":%s}`
const setCredentialResponseJSONWithMetadata = `{"type":"%s","id":"` + uuid + `","name":"%s","value":%s,"metadata":%s,"version_created_at":"` + timestamp + `"}`

func setupSetServer(name, keyType, value string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(fmt.Sprintf(setCredentialRequestJSON, keyType, name, value)),
			RespondWith(http.StatusOK, fmt.Sprintf(setCredentialResponseJSON, keyType, name, value)),
		),
	)
}

func setupSetServerWithMetadata(name, keyType, value, metadata string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(fmt.Sprintf(setCredentialRequestJSONWithMetadata, keyType, name, value, metadata)),
			RespondWith(http.StatusOK, fmt.Sprintf(setCredentialResponseJSONWithMetadata, keyType, name, value, metadata)),
		),
	)
}

func setupPutBadRequestServer(body string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(body),
			RespondWith(http.StatusBadRequest, `{"error":"test error"}`),
		),
	)
}

func testSetFileFailure(caFilename, certificateFilename, privateFilename string) {
	tempDir := test.CreateTempDir("certFilesForTesting")
	if caFilename == "unreadable.txt" {
		caFilename = test.CreateCredentialFile(tempDir, caFilename, "my-ca")
		err := os.Chmod(caFilename, 0222)
		Expect(err).To(BeNil())
	}
	if certificateFilename == "unreadable.txt" {
		certificateFilename = test.CreateCredentialFile(tempDir, certificateFilename, "my-cert")
		err := os.Chmod(certificateFilename, 0222)
		Expect(err).To(BeNil())
	}
	if privateFilename == "unreadable.txt" {
		privateFilename = test.CreateCredentialFile(tempDir, privateFilename, "my-priv")
		err := os.Chmod(privateFilename, 0222)
		Expect(err).To(BeNil())
	}

	session := runCommand("set", "-n", "my-secret",
		"-t", "certificate", "--root", caFilename,
		"--certificate", certificateFilename, "--private", privateFilename)

	os.RemoveAll(tempDir)
	Eventually(session).Should(Exit(1))
	Eventually(session.Err).Should(Say("A referenced file could not be opened. Please validate the provided filenames and permissions, then retry your request."))
}
